-- name: InsertOrder :one
INSERT INTO orders (user_id,status)
VALUES (
           @user_id, @status
       )
RETURNING id;

-- name: InsertItems :one
INSERT INTO items (sku, count)
VALUES (
           @sku, @count
       )
RETURNING id;

-- name: InsertOrderItems :exec
INSERT INTO order_items(order_id, item_id)
VALUES (
           @order_id,@item_id
       );

-- name: UpdateStatusOrders :exec
UPDATE orders SET status = @status WHERE id= @order_id;

-- name: GetInfoFromOrders :one
SELECT
    o.user_id,
    o.status,
    JSON_AGG(JSON_BUILD_OBJECT('sku', i.sku, 'count', i.count)) AS items
FROM orders o
         JOIN order_items oi ON o.id = oi.order_id
         JOIN items i ON oi.item_id = i.id
WHERE o.id= @order_id
GROUP BY o.id;


-- name: ReserveStocks :exec
UPDATE stocks SET reserved= @reserved WHERE id= @sku;

-- name: ReserveRemoveStocks :exec
UPDATE stocks
SET reserved=  @reserved,
    total_count= @total_count
WHERE id= @sku;

-- name: ReserveCancelStocks :exec
UPDATE stocks
SET reserved= @reserved
WHERE id= @sku;

-- name: GetBySKIStocks :one
SELECT total_count, reserved FROM stocks
WHERE id = @sku;
