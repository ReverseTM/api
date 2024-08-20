box.cfg{}

box.once("schema_init", function()
    local space = box.schema.space.create('data_storage', {
        if_not_exists = true,
    })

    -- Создание индекса
    space:create_index('pk', {
        type = 'hash',
        parts = {1, 'string'},
        if_not_exists = true,
    })
end)