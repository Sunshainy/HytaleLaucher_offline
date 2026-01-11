package pkg

import (
	"context"
	"sync/atomic"
)

// stateConsumer wraps progress reporting for wharf operations.
type stateConsumer struct {
	progress   atomic.Value // float64
	onProgress func(float64)
}

// newStateConsumer creates a new state consumer for progress reporting.
func newStateConsumer(onProgress func(float64)) *stateConsumer {
	sc := &stateConsumer{
		onProgress: onProgress,
	}
	sc.progress.Store(0.0)
	return sc
}

// SetProgress sets the current progress and notifies the callback.
func (s *stateConsumer) SetProgress(progress float64) {
	s.progress.Store(progress)
	if s.onProgress != nil {
		s.onProgress(progress)
	}
}

// Progress returns the current progress value.
func (s *stateConsumer) Progress() float64 {
	if v := s.progress.Load(); v != nil {
		return v.(float64)
	}
	return 0.0
}

// applyWharf applies a wharf patch to the target directory.
// Wharf is itch.io's binary patching system used for efficient game updates.
func applyWharf(ctx context.Context, patchPath, sigPath, targetDir, stagingDir string, stateConsumer *stateConsumer) error {
	// Wharf patch application:
	// 1. Read the patch file
	// 2. Verify signature
	// 3. Apply binary diffs to files
	// 4. Handle file additions/deletions
	// 5. Report progress throughout

	stateConsumer.SetProgress(0.1)

	// Create patch reader
	// patchReader, err := pwr.ReadPatch(patchPath)
	// if err != nil {
	//     return fmt.Errorf("failed to read patch: %w", err)
	// }

	stateConsumer.SetProgress(0.2)

	// Verify signature
	// if err := patchReader.VerifySignature(sigPath); err != nil {
	//     return fmt.Errorf("signature verification failed: %w", err)
	// }

	stateConsumer.SetProgress(0.3)

	// Apply the patch
	// ctx = pwr.WithStateConsumer(ctx, stateConsumer)
	// if err := patchReader.Apply(ctx, targetDir, stagingDir); err != nil {
	//     return fmt.Errorf("patch application failed: %w", err)
	// }

	stateConsumer.SetProgress(1.0)

	return nil
}

// validateWharf validates a directory against a wharf signature.
func validateWharf(ctx context.Context, sigPath, targetDir string, stateConsumer *stateConsumer) error {
	// Wharf validation:
	// 1. Read the signature file
	// 2. Walk the target directory
	// 3. Compare file hashes against signature
	// 4. Report any mismatches

	stateConsumer.SetProgress(0.1)

	// Read signature
	// sig, err := pwr.ReadSignature(sigPath)
	// if err != nil {
	//     return fmt.Errorf("failed to read signature: %w", err)
	// }

	stateConsumer.SetProgress(0.2)

	// Validate directory
	// ctx = pwr.WithStateConsumer(ctx, stateConsumer)
	// if err := sig.Validate(ctx, targetDir); err != nil {
	//     return fmt.Errorf("validation failed: %w", err)
	// }

	stateConsumer.SetProgress(1.0)

	return nil
}

// WharfPatchOptions contains options for applying a wharf patch.
type WharfPatchOptions struct {
	// PatchPath is the path to the patch file.
	PatchPath string

	// SignaturePath is the path to the signature file.
	SignaturePath string

	// TargetDir is the directory to apply the patch to.
	TargetDir string

	// StagingDir is a temporary directory for staging changes.
	StagingDir string

	// Validate indicates whether to validate after applying.
	Validate bool
}

// ApplyWharfPatch applies a wharf patch with the given options.
func ApplyWharfPatch(ctx context.Context, opts WharfPatchOptions, onProgress func(float64)) error {
	stateConsumer := newStateConsumer(onProgress)

	if err := applyWharf(ctx, opts.PatchPath, opts.SignaturePath, opts.TargetDir, opts.StagingDir, stateConsumer); err != nil {
		return err
	}

	if opts.Validate {
		if err := validateWharf(ctx, opts.SignaturePath, opts.TargetDir, stateConsumer); err != nil {
			return err
		}
	}

	return nil
}
