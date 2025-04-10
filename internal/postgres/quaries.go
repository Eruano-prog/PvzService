package postgres

import _ "embed"

// User queries
//
//go:embed sql/user_insert.sql
var userInsertQuery string

//go:embed sql/user_find_by_email.sql
var userFindByEmailQuery string

// PVZ queries
//
//go:embed sql/pvz_insert.sql
var pvzInsertQuery string

//go:embed sql/pvz_get_filtered_by_receptionTime_paged.sql
var pvzGetFilteredByReceptionTime string

// Reception queries
//
//go:embed sql/reception_insert.sql
var receptionInsertQuery string

//go:embed sql/reception_get_filtered_by_receptionTime_and_pvzID.sql
var receptionFindByPvzAndTimeQuery string

//go:embed sql/reception_get_active_by_PVZ.sql
var receptionFindByStatusPvzQuery string

//go:embed sql/reception_update_status_by_pvz.sql
var receptionUpdateByStatusAndPvzQuery string

// Product queries
//
//go:embed sql/product_insert_if_reception_active.sql
var productInsertQuery string

//go:embed sql/product_get_by_receptionID.sql
var productFindByReceptionIDQuery string

//go:embed sql/product_delete_if_reception_active.sql
var productDeleteByPvzIDQuery string
