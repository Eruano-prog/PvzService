package postgres

import _ "embed"

// User queries
//
//go:embed sql/user_insert.sql
var UserInsertQuery string

//go:embed sql/user_find_by_email.sql
var UserFindByEmailQuery string

// PVZ queries
//
//go:embed sql/pvz_insert.sql
var PvzInsertQuery string

//go:embed sql/pvz_get_filtered_by_receptionTime_paged.sql
var PvzGetFilteredByReceptionTime string

//go:embed sql/pvz_get_all.sql
var PvzGetAll string

// Reception queries
//
//go:embed sql/reception_insert.sql
var ReceptionInsertQuery string

//go:embed sql/reception_get_filtered_by_receptionTime_and_pvzID.sql
var ReceptionFindByPvzAndTimeQuery string

//go:embed sql/reception_get_active_by_PVZ.sql
var ReceptionFindByStatusPvzQuery string

//go:embed sql/reception_update_status_by_pvz.sql
var ReceptionUpdateByStatusAndPvzQuery string

// Product queries
//
//go:embed sql/product_insert_if_reception_active.sql
var ProductInsertQuery string

//go:embed sql/product_get_by_receptionID.sql
var ProductFindByReceptionIDQuery string

//go:embed sql/product_delete_if_reception_active.sql
var ProductDeleteByPvzIDQuery string
