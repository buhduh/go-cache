package cache

import (
	"testing"
	"time"
)

func TestReaper(t *testing.T) {
	t.Run("MetadataHelper", testHelperFuncs)
	t.Run("Counter", testCounting)
}

func testHelperFuncs(t *testing.T) {
	reaper := newReaper(&NopInvalidator{})
	mData := new(Metadata)
	reaper.Create(mData)
	if mData.Created != time.Now().Unix() {
		t.Errorf(
			"metadata.Created expected %d, got %d",
			time.Now().Unix(), mData.Created,
		)
	}
	if *mData.KeyCount != 1 {
		t.Errorf("wrong count, expected %d got %d", 1, *mData.KeyCount)
	}
	oneSec, _ := time.ParseDuration("1s")
	time.Sleep(oneSec)
	reaper.Update(mData)
	if mData.Modified != time.Now().Unix() {
		t.Errorf(
			"metadata.Updated expected %d, got %d",
			time.Now().Unix(), mData.Modified,
		)
	}
	if *mData.KeyCount != 1 {
		t.Errorf("wrong count, expected %d got %d", 1, *mData.KeyCount)
	}
	time.Sleep(oneSec)
	reaper.Access(mData)
	if mData.Accessed != time.Now().Unix() {
		t.Errorf(
			"metadata.Accessed expected %d, got %d",
			time.Now().Unix(), mData.Accessed,
		)
	}
	if *mData.KeyCount != 1 {
		t.Errorf("wrong count, expected %d got %d", 1, *mData.KeyCount)
	}
}

func testCounting(t *testing.T) {
	reaper := newReaper(&NopInvalidator{})
	mData := new(Metadata)
	for i := 0; i < 200; i++ {
		reaper.Create(mData)
	}
	halfSec, _ := time.ParseDuration(".5s")
	time.Sleep(halfSec)
	if *mData.KeyCount != 200 {
		t.Errorf("count: expected %d, got %d", 200, *mData.KeyCount)
	}
	for i := 0; i < 100; i++ {
		reaper.Remove()
	}
	time.Sleep(halfSec)
	if *mData.KeyCount != 100 {
		t.Errorf("count: expected %d, got %d", 100, *mData.KeyCount)
	}
	reaper.Clear()
	time.Sleep(halfSec)
	if *mData.KeyCount != 0 {
		t.Errorf("count: expected %d, got %d", 0, *mData.KeyCount)
	}
}
