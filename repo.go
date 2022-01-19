package main

import "gopkg.in/mgo.v2/bson"

type repoMemory struct {
	m []contact
}

func newRepoMemory() *repoMemory {
	r := repoMemory{}
	r.m = make([]contact, 0)

	return &r
}

func (r *repoMemory) len() int {
	return len(r.m)
}

func (r *repoMemory) getAll() []contact {
	return r.m[:20]
}

func (r *repoMemory) getByID(id string) contact {
	_, c := r.locate(id)
	return c
}

func (r *repoMemory) store(c *contact) {
	if c.ID == "" {
		c.ID = bson.NewObjectId()
		r.m = append(r.m, *c)
	} else {
		p, c1 := r.locate(c.ID.String())
		r.m[p] = c1
	}
}

func (r *repoMemory) remove(c contact) {
	// todo: not used yet
	//panic(err)
}

func (r *repoMemory) locate(id string) (int, contact) {
	for n, c := range r.m {
		if c.ID.String() == id {
			return n, c
		}
	}
	return -1, contact{}
}

func (r *repoMemory) close() {
	r.m = make([]contact, 0)
}
