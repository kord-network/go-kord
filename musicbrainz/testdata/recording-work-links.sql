COPY(
  SELECT row_to_json(link) FROM (
    SELECT
      'musicbrainz' AS "@source",
      iswc.iswc,
      isrc.isrc
    FROM l_recording_work
    INNER JOIN isrc ON isrc.recording = l_recording_work.entity0
    INNER JOIN iswc ON iswc.work = l_recording_work.entity0
    WHERE l_recording_work.entity1 IN (
      SELECT entity0
      FROM l_recording_work
      INNER JOIN isrc ON isrc.recording = l_recording_work.entity0
      INNER JOIN iswc ON iswc.work = l_recording_work.entity0
      GROUP BY entity0
      HAVING COUNT(entity1) = 3
      LIMIT 5
    )
  ) AS link
) TO STDOUT;
