package crawler

import (
	"github.com/l-dandelion/cwgo/spider"
	"sync"
)

var (
	crawler Crawler
	once    sync.Once
)

type Crawler interface {
	GetSpider(name string) spider.Spider
	AddSpider(sp spider.Spider) error
	DeleteSpider(name string) error
	InitSpider(name string) error
	StartSpider(name string) error
	StopSpider(name string) error
	PauseSpider(name string) error
	RecoverSpider(name string) error
}

type myCrawler struct {
	spiderMapLock sync.RWMutex
	spiderMap     map[string]spider.Spider
}

/*
 * 获取crawler 如果未初始化，先初始化
 */
func New() Crawler {
	once.Do(func() {
		crawler = &myCrawler{
			spiderMap: map[string]spider.Spider{},
		}
	})
	return crawler
}

/*
 * 根据爬虫名获取爬虫
 */
func (crawler *myCrawler) GetSpider(name string) spider.Spider {
	crawler.spiderMapLock.RLock()
	defer crawler.spiderMapLock.RUnlock()
	return crawler.spiderMap[name]
}

/*
 * 添加爬虫
 */
func (crawler *myCrawler) AddSpider(sp spider.Spider) error {
	crawler.spiderMapLock.Lock()
	defer crawler.spiderMapLock.Unlock()
	sp, ok := crawler.spiderMap[sp.Name()]
	if ok {
		return ERR_SPIDER_NAME_REPEATED
	}
	crawler.spiderMap[sp.Name()] = sp
	return nil
}

/*
 * 删除爬虫
 */
func (crawler *myCrawler) DeleteSpider(name string) error {
	crawler.spiderMapLock.Lock()
	defer crawler.spiderMapLock.Unlock()
	_, ok := crawler.spiderMap[name]
	if !ok {
		return ERR_SPIDER_NOT_FOUND
	}
	delete(crawler.spiderMap, name)
	return nil
}

/*
 * 初始化爬虫
 */
func (crawler *myCrawler) InitSpider(name string) error {
	sp := crawler.GetSpider(name)
	if sp == nil {
		return ERR_SPIDER_NOT_FOUND
	}
	return sp.Init()
}

/*
 * 启动爬虫
 */
func (crawler *myCrawler) StartSpider(name string) error {
	sp := crawler.GetSpider(name)
	if sp == nil {
		return ERR_SPIDER_NOT_FOUND
	}
	return sp.Start()
}

/*
 * 终止爬虫
 */
func (crawler *myCrawler) StopSpider(name string) error {
	sp := crawler.GetSpider(name)
	if sp == nil {
		return ERR_SPIDER_NOT_FOUND
	}
	return sp.Stop()
}

/*
 * 暂停爬虫
 */
func (crawler *myCrawler) PauseSpider(name string) error {
	sp := crawler.GetSpider(name)
	if sp == nil {
		return ERR_SPIDER_NOT_FOUND
	}
	return sp.Stop()
}

/*
 * 恢复爬虫
 */
func (crawler *myCrawler) RecoverSpider(name string) error {
	sp := crawler.GetSpider(name)
	if sp == nil {
		return ERR_SPIDER_NOT_FOUND
	}
	return sp.Recover()
}
