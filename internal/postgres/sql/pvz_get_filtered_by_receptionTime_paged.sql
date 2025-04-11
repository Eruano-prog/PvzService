SELECT
    p.*,
    jsonb_agg(
        jsonb_build_object(
            'reception', r,
            'products', (
                SELECT jsonb_agg(pr)
                FROM productions pr
                WHERE pr.reception_id = r.id
            )
        )
    ) AS receptions_with_products
FROM pvz p
INNER JOIN receptions r ON p.id = r.pvz_id
WHERE r.datetime BETWEEN :startTime AND :endTime
GROUP BY p.id
HAVING COUNT(r.id) > 0
OFFSET :offset
LIMIT :limit