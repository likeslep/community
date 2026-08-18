package kafkax

// 事件类型常量（plan.md §18），遵循 `domain.verb` 命名约定。
const (
	EventUserCreated = "user.created"

	EventArticleCreated   = "article.created"
	EventArticleUpdated   = "article.updated"
	EventArticleSubmitted = "article.submitted"
	EventArticlePublished = "article.published"
	EventArticleRejected  = "article.rejected"

	EventQuestionCreated   = "question.created"
	EventQuestionPublished = "question.published"

	EventAnswerCreated  = "answer.created"
	EventAnswerAccepted = "answer.accepted"

	EventCommentCreated    = "comment.created"
	EventLikeCreated       = "like.created"
	EventCollectionCreated = "collection.created"

	EventUserFollowed = "user.followed"
	EventTagFollowed  = "tag.followed"

	EventModerationApproved = "moderation.approved"
	EventModerationRejected = "moderation.rejected"
)
