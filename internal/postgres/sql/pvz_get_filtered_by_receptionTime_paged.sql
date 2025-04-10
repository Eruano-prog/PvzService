SELECT p.*
FROM pvz as p
    inner join receptions as r on p.id = r.pvz_id
WHERE r.date_time BETWEEN :startTime AND :endTime
GROUP BY p.id
HAVING COUNT(r.id) > 0
OFFSET :offset
LIMIT :limit