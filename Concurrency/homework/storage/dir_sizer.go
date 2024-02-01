package storage

import (
	"context"
	"sync"
)

// Result represents the Size function result
type Result struct {
	// Total Size of File objects
	Size int64
	// Count is a count of File objects processed
	Count int64
}

type DirSizer interface {
	// Size calculate a size of given Dir, receive a ctx and the root Dir instance
	// will return Result or error if happened
	Size(ctx context.Context, d Dir) (Result, error)
}

// sizer implement the DirSizer interface
type sizer struct {
	maxWorkersCount int
	sem             chan any
	wg              sync.WaitGroup
	res             chan Result
	err             chan error
}

// NewSizer returns new DirSizer instance
func NewSizer() DirSizer {
	workers := 8
	return &sizer{workers, make(chan any, workers), sync.WaitGroup{},
		make(chan Result, workers), make(chan error)}
}

// Size calculates the size of the given Dir
func (s *sizer) Size(ctx context.Context, d Dir) (Result, error) {

	s.wg.Add(1)
	go func() {
		s.sem <- struct{}{}
		s.worker(ctx, d)
		defer func() {
			<-s.sem
			s.wg.Done()
		}()
	}()

	go func() {
		s.wg.Wait()
		close(s.res)
		close(s.err)
		close(s.sem)
	}()

	return s.listenForResults()
}

// listenForResults combines all the results, listens for errors
func (s *sizer) listenForResults() (Result, error) {
	var total Result
	for {
		closed := false
		select {
		case result, ok := <-s.res:
			if ok {
				total.Size += result.Size
				total.Count += result.Count
			} else {
				closed = true
			}
		case err, ok := <-s.err:
			if ok {
				return Result{}, err
			} else {
				closed = true
			}
		}
		if closed {
			break
		}
	}
	return total, nil
}

// worker is a recursive function that processes the Dir and its subDirs
func (s *sizer) worker(ctx context.Context, dir Dir) {

	dirs := s.processDir(ctx, dir)

	for _, subDir := range dirs {
		select {
		case s.sem <- struct{}{}:
			s.wg.Add(1)
			go func(subDir Dir) {
				s.sem <- struct{}{}
				s.worker(ctx, subDir)
				defer func() {
					<-s.sem
					s.wg.Done()
				}()
			}(subDir)
		default:
			s.worker(ctx, subDir)
		}
	}
}

// processDir calculates the size and the number of files in individual directory
func (s *sizer) processDir(ctx context.Context, dir Dir) []Dir {
	dirs, files, err := dir.Ls(ctx)
	if err != nil {
		s.err <- err
		return nil
	}

	var dirSize int64
	var fileCount int64
	for _, file := range files {
		fileSize, err := file.Stat(ctx)
		if err != nil {
			s.err <- err
			return nil
		}
		dirSize += fileSize
		fileCount++
	}
	s.res <- Result{Size: dirSize, Count: fileCount}

	return dirs
}
