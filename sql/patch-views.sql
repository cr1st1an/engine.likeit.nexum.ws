create view ig_media_likes as select id_ig_media, count(1) as c from likes group by id_ig_media;
