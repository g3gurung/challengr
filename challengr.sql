/*
(((SELECT COUNT(posts.id) FROM posts WHERE posts.challenge_id=id AND (SELECT COUNT(likes.id) FROM likes WHERE likes.post_id=posts.id)) / (SELECT COUNT(posts.id) FROM posts WHERE posts.challenge_id=id)) * (SELECT COUNT(posts.id) FROM posts WHERE posts.challenge_id=id) / 100) + (((SELECT COUNT(posts.id) FROM posts WHERE posts.challenge_id=id AND (SELECT COUNT(likes.id) FROM likes WHERE likes.post_id=posts.id)) / (SELECT COUNT(posts.id) FROM posts WHERE posts.challenge_id=id)) * (SELECT COUNT(posts.id) FROM posts WHERE posts.challenge_id=id) / 100)*challenge.wieght
*/


