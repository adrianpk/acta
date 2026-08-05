package acta

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

const maxEventLineSize = 8 << 20

var (
	// ErrJournalNotLoaded indicates that journal IO started before Load.
	ErrJournalNotLoaded = errors.New("journal is not loaded")
	// ErrJournalCorrupt indicates invalid canonical journal content.
	ErrJournalCorrupt = errors.New("journal is corrupt")
)

// MilestoneKey is the durable milestone identity.
type MilestoneKey struct {
	RepositoryIdentity string
	Commit             string
}

// AppendResult reports whether Append published a new event.
type AppendResult struct {
	Event    Event
	Recorded bool
}

// Journal owns monthly canonical segments and disposable indexes.
type Journal struct {
	config Config

	mu         sync.RWMutex
	loaded     bool
	events     map[EventID]Event
	milestones map[MilestoneKey]EventID
	effective  map[MilestoneKey]EventID
	order      []EventID
}

// NewJournal stores configuration without performing IO.
func NewJournal(config Config) *Journal {
	return &Journal{config: config}
}

// Load creates Acta-owned directories and rebuilds all journal indexes.
func (journal *Journal) Load(ctx context.Context) error {
	journal.mu.Lock()
	defer journal.mu.Unlock()

	journal.loaded = false
	journal.applyIndex(newJournalIndex())

	if err := ctx.Err(); err != nil {
		return err
	}

	if err := prepareJournalDirectories(journal.config); err != nil {
		return err
	}

	if _, err := recoverPending(journal.config); err != nil {
		return err
	}

	index, err := rebuildJournalIndex(ctx, journal.config.JournalRoot)
	if err != nil {
		return err
	}

	journal.applyIndex(index)

	journal.loaded = true

	return nil
}

// Append validates, publishes, syncs, and indexes one event.
func (journal *Journal) Append(ctx context.Context, event Event) (AppendResult, error) {
	journal.mu.Lock()
	defer journal.mu.Unlock()

	if !journal.loaded {
		return AppendResult{}, ErrJournalNotLoaded
	}

	if err := journal.recoverForAppend(ctx); err != nil {
		return AppendResult{}, err
	}

	if err := ctx.Err(); err != nil {
		return AppendResult{}, err
	}

	if err := event.Validate(); err != nil {
		return AppendResult{}, err
	}

	if existing, ok := journal.existingMilestone(event); ok {
		return AppendResult{Event: existing, Recorded: false}, nil
	}

	if _, ok := journal.events[event.ID]; ok {
		return AppendResult{}, fmt.Errorf("append event: duplicate event ID %q", event.ID)
	}

	if err := validateEventAgainstIndex(event, journal.events, journal.milestones, journal.effective); err != nil {
		return AppendResult{}, err
	}

	line, err := json.Marshal(event)
	if err != nil {
		return AppendResult{}, fmt.Errorf("encode event: %w", err)
	}

	line = append(line, '\n')

	segmentPath := SegmentPath(journal.config.JournalRoot, event.RecordedAt)
	if err := publishEvent(journal.config, segmentPath, line); err != nil {
		return AppendResult{}, err
	}

	addEventToIndex(event, journal.events, journal.milestones, journal.effective, &journal.order)

	return AppendResult{Event: event, Recorded: true}, nil
}

func (journal *Journal) existingMilestone(event Event) (Event, bool) {
	if event.Type != EventTypeMilestone {
		return Event{}, false
	}

	existingID, ok := journal.milestones[eventMilestoneKey(event)]
	if !ok {
		return Event{}, false
	}

	return journal.events[existingID], true
}

func (journal *Journal) recoverForAppend(ctx context.Context) error {
	recovered, err := recoverPending(journal.config)
	if err != nil || !recovered {
		return err
	}

	index, err := rebuildJournalIndex(ctx, journal.config.JournalRoot)
	if err != nil {
		return err
	}

	journal.applyIndex(index)

	return nil
}

// Event returns one event by ID.
func (journal *Journal) Event(id EventID) (Event, bool) {
	journal.mu.RLock()
	defer journal.mu.RUnlock()

	event, ok := journal.events[id]

	return event, ok
}

// Milestone returns the original milestone for a repository and commit.
func (journal *Journal) Milestone(key MilestoneKey) (Event, bool) {
	journal.mu.RLock()
	defer journal.mu.RUnlock()

	id, ok := journal.milestones[key]
	if !ok {
		return Event{}, false
	}

	return journal.events[id], true
}

// Effective returns the current effective event for a milestone.
func (journal *Journal) Effective(key MilestoneKey) (Event, bool) {
	journal.mu.RLock()
	defer journal.mu.RUnlock()

	id, ok := journal.effective[key]
	if !ok {
		return Event{}, false
	}

	return journal.events[id], true
}

// Count returns the number of canonical events loaded in memory.
func (journal *Journal) Count() int {
	journal.mu.RLock()
	defer journal.mu.RUnlock()

	return len(journal.events)
}

// SegmentPath returns the canonical monthly path for a recording time.
func SegmentPath(root string, recordedAt time.Time) string {
	utc := recordedAt.UTC()

	return filepath.Join(root, fmt.Sprintf("%04d", utc.Year()), fmt.Sprintf("%02d.jsonl", int(utc.Month())))
}

type journalIndex struct {
	events     map[EventID]Event
	milestones map[MilestoneKey]EventID
	effective  map[MilestoneKey]EventID
	order      []EventID
}

func newJournalIndex() *journalIndex {
	return &journalIndex{
		events:     make(map[EventID]Event),
		milestones: make(map[MilestoneKey]EventID),
		effective:  make(map[MilestoneKey]EventID),
		order:      make([]EventID, 0),
	}
}

func rebuildJournalIndex(ctx context.Context, root string) (*journalIndex, error) {
	paths, err := listSegmentPaths(root)
	if err != nil {
		return nil, err
	}

	index := newJournalIndex()
	for _, segmentPath := range paths {
		if err := loadSegment(ctx, root, segmentPath, index); err != nil {
			return nil, err
		}
	}

	return index, nil
}

func (journal *Journal) applyIndex(index *journalIndex) {
	journal.events = index.events
	journal.milestones = index.milestones
	journal.effective = index.effective
	journal.order = index.order
}

func prepareJournalDirectories(config Config) error {
	for _, directory := range []string{config.FilesRoot, config.JournalRoot} {
		if directory == "" {
			return fmt.Errorf("prepare journal: configured path is empty")
		}

		if err := os.MkdirAll(directory, 0o700); err != nil {
			return fmt.Errorf("prepare journal directory %q: %w", directory, err)
		}

		if err := os.Chmod(directory, 0o700); err != nil {
			return fmt.Errorf("set journal directory permissions %q: %w", directory, err)
		}
	}

	return nil
}

func listSegmentPaths(root string) ([]string, error) {
	var paths []string

	err := filepath.WalkDir(root, func(current string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		if entry.IsDir() {
			return nil
		}

		if strings.HasSuffix(entry.Name(), ".jsonl") {
			paths = append(paths, current)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("list journal segments: %w", err)
	}

	sort.Strings(paths)

	return paths, nil
}

func loadSegment(ctx context.Context, root, segmentPath string, index *journalIndex) (resultErr error) {
	file, err := os.Open(segmentPath)
	if err != nil {
		return fmt.Errorf("open journal segment %q: %w", segmentPath, err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil && resultErr == nil {
			resultErr = fmt.Errorf("close journal segment %q: %w", segmentPath, closeErr)
		}
	}()

	if err := validateSegmentTail(file, segmentPath); err != nil {
		return err
	}

	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 64<<10), maxEventLineSize)

	lineNumber := 0
	for scanner.Scan() {
		lineNumber++

		if err := ctx.Err(); err != nil {
			return err
		}

		if err := loadSegmentEvent(root, segmentPath, lineNumber, scanner.Bytes(), index); err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		return corruptSegment(segmentPath, lineNumber+1, "read event: %v", err)
	}

	return nil
}

func loadSegmentEvent(root, segmentPath string, lineNumber int, line []byte, index *journalIndex) error {
	if len(line) == 0 {
		return corruptSegment(segmentPath, lineNumber, "blank line")
	}

	var event Event
	if err := json.Unmarshal(line, &event); err != nil {
		return corruptSegment(segmentPath, lineNumber, "decode event: %v", err)
	}

	if err := event.Validate(); err != nil {
		return corruptSegment(segmentPath, lineNumber, "%v", err)
	}

	if SegmentPath(root, event.RecordedAt) != segmentPath {
		return corruptSegment(segmentPath, lineNumber, "event belongs to %q", SegmentPath(root, event.RecordedAt))
	}

	if _, ok := index.events[event.ID]; ok {
		return corruptSegment(segmentPath, lineNumber, "duplicate event ID %q", event.ID)
	}

	if err := validateEventAgainstIndex(event, index.events, index.milestones, index.effective); err != nil {
		return corruptSegment(segmentPath, lineNumber, "%v", err)
	}

	addEventToIndex(event, index.events, index.milestones, index.effective, &index.order)

	return nil
}

func validateEventAgainstIndex(event Event, events map[EventID]Event, milestones, effective map[MilestoneKey]EventID) error {
	key := eventMilestoneKey(event)
	switch event.Type {
	case EventTypeMilestone:
		if _, ok := milestones[key]; ok {
			return fmt.Errorf("duplicate milestone key for repository %q and commit %q", key.RepositoryIdentity, key.Commit)
		}
	case EventTypeCorrection:
		target, ok := events[event.TargetID]
		if !ok {
			return fmt.Errorf("correction target %q does not exist", event.TargetID)
		}

		targetKey := eventMilestoneKey(target)
		if key != targetKey {
			return fmt.Errorf("correction evidence does not match target %q", event.TargetID)
		}

		if effective[targetKey] != event.TargetID {
			return fmt.Errorf("correction target %q is not the current effective event", event.TargetID)
		}

		if !sameCorrectionSnapshot(event, target) {
			return fmt.Errorf("correction snapshot does not match target %q", event.TargetID)
		}
	}

	return nil
}

func addEventToIndex(event Event, events map[EventID]Event, milestones, effective map[MilestoneKey]EventID, order *[]EventID) {
	events[event.ID] = event
	*order = append(*order, event.ID)

	key := eventMilestoneKey(event)
	if event.Type == EventTypeMilestone {
		milestones[key] = event.ID
	}

	effective[key] = event.ID
}

func eventMilestoneKey(event Event) MilestoneKey {
	return MilestoneKey{RepositoryIdentity: event.Repository.Identity, Commit: event.Git.Commit}
}

func sameCorrectionSnapshot(correction, target Event) bool {
	if correction.Repository != target.Repository || correction.Git != target.Git {
		return false
	}

	if correction.Execution == nil || target.Execution == nil {
		return correction.Execution == nil && target.Execution == nil
	}

	return *correction.Execution == *target.Execution
}

func appendEventLine(segmentPath string, line []byte) error {
	if err := os.MkdirAll(filepath.Dir(segmentPath), 0o700); err != nil {
		return fmt.Errorf("create segment directory: %w", err)
	}

	if err := os.Chmod(filepath.Dir(segmentPath), 0o700); err != nil {
		return fmt.Errorf("set segment directory permissions: %w", err)
	}

	file, err := os.OpenFile(segmentPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600)
	if err != nil {
		return fmt.Errorf("open journal segment %q: %w", segmentPath, err)
	}

	if err := file.Chmod(0o600); err != nil {
		return closeJournalFile(file, fmt.Errorf("set journal segment permissions %q: %w", segmentPath, err))
	}

	if _, err := file.Write(line); err != nil {
		return closeJournalFile(file, fmt.Errorf("append journal segment %q: %w", segmentPath, err))
	}

	if err := file.Sync(); err != nil {
		return closeJournalFile(file, fmt.Errorf("sync journal segment %q: %w", segmentPath, err))
	}

	return closeJournalFile(file, nil)
}

func closeJournalFile(file *os.File, operationErr error) error {
	closeErr := file.Close()

	return errors.Join(operationErr, closeErr)
}

func corruptSegment(segmentPath string, line int, format string, arguments ...any) error {
	message := fmt.Sprintf(format, arguments...)

	return fmt.Errorf("%w: %s:%d: %s", ErrJournalCorrupt, segmentPath, line, message)
}
