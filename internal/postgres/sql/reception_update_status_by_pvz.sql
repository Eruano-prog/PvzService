UPDATE receptions
SET status = :status
WHERE pvz_id = :pvz_id
    AND status = :last_status
RETURNING *;