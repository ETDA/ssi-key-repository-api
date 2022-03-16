import * as Knex from "knex";


export async function up(knex: Knex): Promise<void> {
    return knex.schema.createTable("keys", function (table) {
        table.string('id', 255).primary()
        table.text('public_key').notNullable()
        table.text('private_key_encrypted').notNullable()
        table.string('type', 255).notNullable()
        table.dateTime('created_at').notNullable()
        table.dateTime('updated_at').notNullable()
        table.dateTime('deleted_at')
    })
}


export async function down(knex: Knex): Promise<void> {
    return knex.schema.dropTableIfExists('keys')
}

