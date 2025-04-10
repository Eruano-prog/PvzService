SELECT *
FROM receptions as r
WHERE pvz_id = :pvz_id
  AND datetime BETWEEN :from_time and :to_time