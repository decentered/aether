// Services > Configstore > Content Relations

// This package is the content equivalent of user relations, it handles actions like subscribing and unsubscribing to a board, following a thread, and so on.

// This package also handles the optional whitelist for boards, which is used to highlight especially interesting communities. Neither this whitelist, nor the API for it is *not* a part of the protocol, only a part of this specific c0 client app.

/**
 *
 * Heads up - the way to use this is through configstore and you have to access this through GetContentRelations, and when done, you should do SetContentRelations, otherwise it won't be saved permanently.
 *
 * If you want to just read, you can read without doing a Get. But if you want to read and write, you should do get, edit and set, because that's the only way to retain changes.
 *
 */

package configstore

import (
	"sync"
)

type Board struct {
	Fingerprint string
	Notify      bool
	LastSeen    int64
}

type Thread struct {
	Fingerprint string
	Notify      bool
	LastSeen    int64
}

type ContentRelations struct {
	lock          sync.Mutex
	Initialised   bool
	SubbedBoards  []Board
	SubbedThreads []Thread
	Whitelist     whitelist
}

type whitelist struct {
	Disabled   bool
	LastUpdate int64
	Source     string
	Boards     []string
}

func (c *ContentRelations) Init() {
	c.Initialised = true
}

func (wl *whitelist) Update() {
	if wl.Disabled {
		return
	}
	if len(wl.Source) == 0 {
		wl.Source = "www.example.com/whitelist" // todo
	}
	// If last update was less than an hour ago, return
	// Check the whitelist source
}

// func (c *ContentRelations) GetAllSubbedBoards() []Board {
// 	c.lock.Lock()
// 	defer c.lock.Unlock()
// 	return c.SubbedBoards
// }

// func (c *ContentRelations) GetAllSubbedThreads() []Thread {
// 	c.lock.Lock()
// 	defer c.lock.Unlock()
// 	return c.SubbedThreads
// }

func (c *ContentRelations) RefreshWhitelist() {
	c.Whitelist.Boards = []string{}
	c.Whitelist.Boards = []string{"00886d50e598e43984d0df37f83b2398d371a9cc8417a9bba521a95c2da45ffe"} // debug
}

func (c *ContentRelations) IsSubbedBoard(fp string) (isSubbed, notifyEnabled bool, lastSeen int64) {
	loc := c.FindBoard(fp)
	if loc != -1 {
		return true, c.SubbedBoards[loc].Notify, c.SubbedBoards[loc].LastSeen
	}
	return false, false, 0
}

func (c *ContentRelations) IsWhitelistedBoard(fp string) (isWhitelisted bool) {
	c.RefreshWhitelist()
	return c.FindBoardInWhitelist(fp) != -1
}

func (c *ContentRelations) FindBoard(fp string) int {
	for key, _ := range c.SubbedBoards {
		if c.SubbedBoards[key].Fingerprint == fp {
			return key
		}
	}
	return -1
}

func (c *ContentRelations) FindBoardInWhitelist(fp string) int {
	for key, _ := range c.Whitelist.Boards {
		if c.Whitelist.Boards[key] == fp {
			return key
		}
	}
	return -1
}

func (c *ContentRelations) FindThread(fp string) int {
	for key, _ := range c.SubbedThreads {
		if c.SubbedThreads[key].Fingerprint == fp {
			return key
		}
	}
	return -1
}

func (c *ContentRelations) insertBoard(fp string, notify bool, lastseen int64, lastSeenOnly bool) {
	if i := c.FindBoard(fp); i != -1 {
		c.SubbedBoards[i].Notify = notify
		if lastseen > c.SubbedBoards[i].LastSeen {
			c.SubbedBoards[i].LastSeen = lastseen
		}
		return
	}
	c.SubbedBoards = append(c.SubbedBoards,
		Board{Fingerprint: fp, Notify: notify})
}

func (c *ContentRelations) insertLastSeenForBoard(fp string, lastseen int64) {
	if i := c.FindBoard(fp); i != -1 {
		if lastseen > c.SubbedBoards[i].LastSeen {
			c.SubbedBoards[i].LastSeen = lastseen
		}
	}
}

func (c *ContentRelations) insertThread(fp string, notify bool) {
	if i := c.FindThread(fp); i != -1 {
		c.SubbedThreads[i].Notify = notify
		return
	}
	c.SubbedThreads = append(c.SubbedThreads,
		Thread{Fingerprint: fp, Notify: notify})
}

func (c *ContentRelations) removeBoard(fp string) {
	if i := c.FindBoard(fp); i != -1 {
		c.SubbedBoards = append(c.SubbedBoards[0:i], c.SubbedBoards[i+1:len(c.SubbedBoards)]...)
	}
}

func (c *ContentRelations) removeThread(fp string) {
	if i := c.FindThread(fp); i != -1 {
		c.SubbedThreads = append(c.SubbedThreads[0:i], c.SubbedThreads[i+1:len(c.SubbedThreads)]...)
	}
}

// SetBoardSignal sets the board signal into the storage. If a board is subscribed, we set the notify signal as well, if a subscription is removed, we remove the entry.
func (c *ContentRelations) SetBoardSignal(
	fp string, subscribed, notify bool, lastseen int64, lastSeenOnly bool) (committed bool) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if lastSeenOnly {
		c.insertLastSeenForBoard(fp, lastseen)
		return
	}
	if subscribed {
		c.insertBoard(fp, notify, lastseen, lastSeenOnly)
	} else {
		c.removeBoard(fp)
	}
	return true
}

// func (c *ContentRelations) SubBoard(fp string, notify bool) {
// 	c.insertBoard(fp, notify)
// }

// func (c *ContentRelations) UnsubBoard(fp string) {
// 	// c.lock.Lock()
// 	// defer c.lock.Unlock()
// 	c.removeBoard(fp)
// }

// func (c *ContentRelations) SubThread(fp string, notify bool) {
// 	c.lock.Lock()
// 	defer c.lock.Unlock()
// 	c.insertThread(fp, notify)
// }

// func (c *ContentRelations) UnsubThread(fp string) {
// 	c.lock.Lock()
// 	defer c.lock.Unlock()
// 	c.removeThread(fp)
// }
