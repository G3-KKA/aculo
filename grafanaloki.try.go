package main


redisRes := r.rdb.Get(r.ctx, "article:"+idStr)

if redisRes == nil {
	redisRes.Scan(article)
	slog.Debug("get article from cache")
	return article, nil
}

// если статьи нету в кэше, то делаем запрос в бд
err := pgxscan.Get(
	r.ctx, r.pg, article, `
	SELECT *
	FROM articles 
	WHERE id=$1;`, id,
)
if err != nil {
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, types.ErrArticleNotFound
	}
	return nil, err
}
// кэшируем полученную из бд статью
err = r.rdb.Set(r.ctx, "article:"+idStr, article, 0).Err()
if err != nil {
	return nil, err
}
return article, nil
