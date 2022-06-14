package article

import (
	"cms/articles"
	repo "cms/articles/db"
	"database/sql"
	"fmt"
)

type articleTag struct {
	Slug string `json:"slug"`
	Name string `json:"name"`
}

type withArticleTags struct {
	tagsRepository *repo.TagRepository
}

func (w withArticleTags) getOrCreateTags(tx *sql.Tx, tags []articleTag) ([]int, error) {
	var ids []int
	uniqIDs := make(map[int]struct{})
	for _, tag := range tags {
		tagDTO := repo.Tag{
			Name: tag.Name,
			Slug: tag.Slug,
		}
		if "" == tagDTO.Slug {
			tagDTO.Slug = articles.GenerateSlug(tagDTO.Name)
		}
		id, err := w.tagsRepository.GetOrCreateTag(tx, tagDTO)
		if nil != err {
			return nil, fmt.Errorf("get or create tag: %s", err)
		}
		if _, ok := uniqIDs[id]; !ok {
			uniqIDs[id] = struct{}{}
			ids = append(ids, id)
		}
	}
	return ids, nil
}

func (w withArticleTags) mergeTags(existed []int, updated []int) (toInsert []int, toDelete []int) {
	existedMap := make(map[int]struct{})
	for _, id := range existed {
		existedMap[id] = struct{}{}
	}
	for _, id := range updated {
		if _, ok := existedMap[id]; ok {
			delete(existedMap, id)
		} else {
			toInsert = append(toInsert, id)
		}
	}
	for id := range existedMap {
		toDelete = append(toDelete, id)
	}
	return
}
