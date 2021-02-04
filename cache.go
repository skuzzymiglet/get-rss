package main

import "github.com/peterbourgon/diskv"

// DiskvCacheFinder caches a child finder's results with diskv
type DiskvCacheFinder struct {
	Child FeedFinder
	Opts  diskv.Options
	cache *diskv.Diskv
}

func (d *DiskvCacheFinder) Init() error {
	d.cache = diskv.New(d.Opts)
	return d.Child.Init()
}

func (d *DiskvCacheFinder) Find(query string) (string, error) {
	if d.cache.Has(query) {
		b, err := d.cache.Read(query)
		return string(b), err
	}
	s, err := d.Child.Find(query)
	if err != nil {
		return s, err
	}
	// TODO: slugify names
	return s, d.cache.Write(query, []byte(s))
}
