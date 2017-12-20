COPY(
  SELECT row_to_json(artist) FROM (
    SELECT
      'musicbrainz' AS "@source",
      artist.id,
      artist.name,
      artist.sort_name,
      artist_type.name AS type,
      gender.name AS gender,
      area.name AS area,
      CASE
	WHEN artist.begin_date_year IS NULL THEN NULL
	ELSE make_date(
	  artist.begin_date_year,
	  CASE
	    WHEN artist.begin_date_month IS NULL THEN 1
	    ELSE artist.begin_date_month
	  END,
	  CASE
	    WHEN artist.begin_date_day IS NULL THEN 1
	    ELSE artist.begin_date_day
	  END
	)
      END AS begin_date,
      CASE
	WHEN artist.end_date_year IS NULL THEN NULL
	ELSE make_date(
	  artist.end_date_year,
	  CASE
	    WHEN artist.end_date_month IS NULL THEN 1
	    ELSE artist.end_date_month
	  END,
	  CASE
	    WHEN artist.end_date_day IS NULL THEN 1
	    ELSE artist.end_date_day
	  END
	)
      END AS end_date,
      ARRAY(SELECT ipi FROM artist_ipi WHERE artist_ipi.artist = artist.id) AS ipi,
      ARRAY(SELECT isni FROM artist_isni WHERE artist_isni.artist = artist.id) AS isni,
      ARRAY(SELECT name FROM artist_alias WHERE artist_alias.artist = artist.id) AS alias,
      artist.gid AS mbid,
      artist.comment AS disambiguation_comment,
      ARRAY(SELECT text FROM annotation LEFT JOIN artist_annotation ON artist_annotation.annotation = annotation.id WHERE artist_annotation.artist = artist.id) AS annotation
    FROM artist
    LEFT JOIN artist_type ON artist.type = artist_type.id
    LEFT JOIN gender ON artist.gender = gender.id
    LEFT JOIN area ON artist.area = area.id
    WHERE artist.id IN (
      389042,
      655,
      800841,
      11035,
      11669,
      32197
    )
  ) AS artist
) TO STDOUT;
