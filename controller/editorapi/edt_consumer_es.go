package editorapi

import (
	"editorApi/init/mgdb"
	"editorApi/requests"
	"strings"
	"tkCommon/lib/els"

	"go.mongodb.org/mongo-driver/bson"
)

func PushContentToEs(msg *requests.DictESParams) {

	indexName := "index_" + msg.From + "2" + msg.To

	switch msg.Operate {
	case "offline":
		fallthrough
	case "delete":
		for _, id := range msg.Uuids {
			els.DelById(indexName, id)
		}
	case "online":

		switch msg.CType {
		case "word":
			words := []map[string]interface{}{}
			mgdb.Find(
				mgdb.EnvEditor,
				mgdb.DbDict,
				"dict_"+msg.From,
				bson.M{
					"uuid": bson.M{
						"$in": msg.Uuids,
					},
				},
				nil,
				nil,
				0,
				1000,
				&words,
			)
			for _, word := range words {
				var body map[string]interface{}
				id := word["uuid"].(string)

				body = map[string]interface{}{
					"id":        id,
					"content":   strings.Trim(word["content"].(string), " "),
					"is_word":   true,
					"ctype":     msg.CType,
					"lang_code": msg.From,
				}

				if cid, ok := word["card_id"]; ok {
					body["card_id"] = cid
				}
				if imgs, ok := word["images"]; ok {
					body["images"] = imgs
				}
				wordAttr := map[string]interface{}{}

				mgdb.FindOne(
					mgdb.EnvEditor,
					mgdb.DbDict,
					"dict_"+msg.From+"_"+msg.To,
					bson.M{
						"parent": id,
					},
					nil,
					&wordAttr,
				)

				body["tags"] = wordAttr["tags"]
				extras := map[string]interface{}{
					"sound_infos": word["sound_infos"],
					"content_trs": wordAttr["content_tr"],
					"synonyms":    wordAttr["synonym"],
					"homonyms":    wordAttr["homonyms"],
				}

				//相关短语
				phrases := []map[string]string{}
				mgdb.Find(
					mgdb.EnvEditor,
					mgdb.DbDict,
					"phrase_"+msg.From,
					bson.M{
						"dict_uuid": id,
					},
					nil,
					map[string]int{
						"content": 1,
						"uuid":    1,
						"_id":     0,
					},
					0,
					1000,
					&phrases,
				)

				phrasesExtras := make([]map[string]string, len(phrases))
				for k, p := range phrases {
					phraseAttr := map[string]string{}
					mgdb.FindOne(
						mgdb.EnvEditor,
						mgdb.DbDict,
						"phrase_"+msg.From+"_"+msg.To,
						bson.M{
							"parent": p["uuid"],
						},
						nil,
						&phraseAttr,
					)
					phrasesExtras[k] = map[string]string{
						"phrase":   p["content"],
						"phraseTr": phraseAttr["content_tr"],
					}
				}
				extras["phrases"] = phrasesExtras
				body["extras"] = extras
				els.Upsert(indexName, id, body)
			}
		case "sentence":
			sentences := []map[string]interface{}{}
			mgdb.Find(
				mgdb.EnvEditor,
				mgdb.DbDict,
				"sentence_"+msg.From,
				bson.M{
					"uuid": bson.M{"$in": msg.Uuids},
				},
				nil,
				nil,
				0,
				100,
				&sentences,
			)
			for _, sentence := range sentences {
				var body map[string]interface{}
				id := sentence["uuid"].(string)
				body = map[string]interface{}{
					"id":          id,
					"card_id":     sentence["card_id"],
					"con_from":    sentence["con_from"],
					"content":     strings.Trim(sentence["sentence"].(string), " "),
					"is_sentence": true,
					"ctype":       "sentence",
					"lang_code":   msg.From,
					"tags":        sentence["tags"],
					"images":      sentence["images"],
				}
				if cid, ok := sentence["card_id"]; ok {
					body["card_id"] = cid
				}
				if imgs, ok := sentence["image"]; ok {
					body["images"] = imgs
				}
				extras := map[string]interface{}{
					"source": sentence["source"],
				}
				if sounds, ok := sentence["sound_infos"]; ok {
					extras["sound_infos"] = sounds
				}
				sentenceAttr := map[string]string{}
				mgdb.FindOne(
					mgdb.EnvEditor,
					mgdb.DbDict,
					"sentence_"+msg.From+"_"+msg.To,
					bson.M{
						"parent": id,
					},
					nil,
					&sentenceAttr,
				)
				extras["content_tr"] = sentenceAttr["content_tr"]
				body["extras"] = extras
				els.Upsert(indexName, id, body)
			}
		}

	}
}
