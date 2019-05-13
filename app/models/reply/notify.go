package reply

import (
	"errors"
	notification "gin_bbs/app/models/notification"
	topicModel "gin_bbs/app/models/topic"
	userModel "gin_bbs/app/models/user"
	"strconv"
)

// TopicRepliedNotify -
func TopicRepliedNotify(reply *Reply, currentUser *userModel.User) error {
	data, _, user, err := getTopicReplied(reply, currentUser)
	if err != nil {
		return err
	}

	if err = notification.Notify("TopicReplied", userModel.TableName, reply.UserID, data); err != nil {
		return err
	}

	// user notification_count ++
	user.NotificationCount = user.NotificationCount + 1
	return user.Update()
}

// ----------------------- private
func getTopicReplied(reply *Reply, currentUser *userModel.User) (map[string]interface{}, *topicModel.Topic, *userModel.User, error) {
	topic, err := topicModel.Get(int(reply.TopicID))
	if err != nil {
		return nil, nil, nil, err
	}
	// 要通知的人是当前用户 (帖子主人正是当前发布 reply 的用户)
	if topic.UserID == currentUser.ID {
		return nil, nil, nil, errors.New("帖子主人正是当前发布评论的用户")
	}

	user, err := userModel.Get(int(reply.UserID))
	if err != nil {
		return nil, nil, nil, err
	}

	data := map[string]interface{}{
		"reply_id":      reply.ID,
		"reply_content": reply.Content,
		"user_id":       reply.UserID,
		"user_name":     user.Name,
		"user_avatar":   user.Avatar,
		"topic_link":    topic.Link() + "#reply" + strconv.Itoa(int(reply.ID)),
		"topic_id":      topic.ID,
		"topic_title":   topic.Title,
	}

	return data, topic, user, nil
}