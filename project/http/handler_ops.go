// package http

// import (
// 	"net/http"

// 	"github.com/labstack/echo/v4"
// )

// func (h Handler) AllOpsBookings(ctx echo.Context) error {
// 	bookings, err := h.opsBookingReadModel.AllBookings()
// 	if err != nil {
// 		return err
// 	}

// 	return ctx.JSON(http.StatusOK, bookings)
// }

// func (h Handler) GetOpsBooking(ctx echo.Context) error {
// 	id := ctx.Param("id")

// 	booking, err := h.opsBookingReadModel.ReservationReadModel(ctx.Request().Context(), id)
// 	if err != nil {
// 		return err
// 	}

//		return ctx.JSON(http.StatusOK, booking)
//	}
package http
