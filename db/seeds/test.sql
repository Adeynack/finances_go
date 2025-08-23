INSERT INTO "public"."users" ("id", "created_at", "updated_at", "email", "encrypted_password", "admin", "display_name") VALUES
('569bcfdd-4056-42cd-af9c-285fa5ce92c8', '2025-08-14 21:46:07.62893', '2025-08-14 21:46:07.62893', 'joe@example.com', 'asd', 'f', 'Joe'),
('5b2fb2a0-bdc1-454d-af50-067e9b9e9dd3', '2025-08-14 21:46:53.577442', '2025-08-14 21:46:53.577442', 'linda@example.com', 'qwe', 'f', 'Linda');

INSERT INTO "public"."books" ("id", "created_at", "updated_at", "name", "owner_id", "default_currency_iso_code") VALUES
('8d8666c0-016f-49fb-8f59-4150a822ffb2', '2025-08-14 21:47:58.211393', '2025-08-14 21:47:58.211393', 'Joe''s Book', '569bcfdd-4056-42cd-af9c-285fa5ce92c8', 'EUR'),
('5da6e20f-eecd-456b-a8dd-ae1a63d0268e', '2025-08-24 00:20:00.000000', '2025-08-24 00:20:00.000000', 'Foo, the Book', '569bcfdd-4056-42cd-af9c-285fa5ce92c8', 'USD');

INSERT INTO "public"."registers" ("id", "created_at", "updated_at", "name", "type", "book_id", "parent_id", "starts_at", "expires_at", "currency_iso_code", "notes", "initial_balance", "active", "default_category", "institution_name", "account_number", "iban", "annual_interest_rate", "credit_limit", "card_number") VALUES
('7625b655-732d-49e0-a86b-43994ea89359', '2025-08-14 21:50:02.935459', '2025-08-14 21:50:02.935459', 'Joe''s Credit Card', 'Card', '8d8666c0-016f-49fb-8f59-4150a822ffb2', NULL, '2025-08-14', NULL, 'EUR', NULL, 0, 't', NULL, NULL, NULL, NULL, NULL, NULL, NULL),
('af903dfd-3e65-41d2-83f8-db1422b5159c', '2025-08-14 21:52:49.46216', '2025-08-14 21:52:49.46216', 'Alcohol', 'Expense', '8d8666c0-016f-49fb-8f59-4150a822ffb2', NULL, '2025-08-14', NULL, 'EUR', NULL, 0, 't', NULL, NULL, NULL, NULL, NULL, NULL, NULL),
('df622be3-dff4-4541-b2c7-40947bc68d5f', '2025-08-14 21:52:49.46216', '2025-08-14 21:52:49.46216', 'Home Improvements', 'Expense', '8d8666c0-016f-49fb-8f59-4150a822ffb2', NULL, '2025-08-14', NULL, 'EUR', NULL, 0, 't', NULL, NULL, NULL, NULL, NULL, NULL, NULL),
('f01ac27a-89f4-4aab-b534-0cf77cee659c', '2025-08-14 21:52:49.46216', '2025-08-14 21:52:49.46216', 'Food', 'Expense', '8d8666c0-016f-49fb-8f59-4150a822ffb2', NULL, '2025-08-14', NULL, 'EUR', NULL, 0, 't', NULL, NULL, NULL, NULL, NULL, NULL, NULL),
('f7baf4ac-52f8-494c-87b3-5178b416411f', '2025-08-14 21:48:56.919668', '2025-08-14 21:48:56.919668', 'Joe''s Main Bank Account', 'Bank', '8d8666c0-016f-49fb-8f59-4150a822ffb2', NULL, '2025-08-14', NULL, 'EUR', NULL, 0, 't', NULL, NULL, NULL, NULL, NULL, NULL, NULL);

INSERT INTO "public"."exchanges" ("id", "created_at", "updated_at", "date", "register_id", "cheque", "description", "memo", "status") VALUES
('4c703f3b-7505-4785-9c63-14b38b9e1129', '2025-08-14 21:53:51.988009', '2025-08-14 21:53:51.988009', '2025-08-14', '7625b655-732d-49e0-a86b-43994ea89359', NULL, 'Grocery', NULL, 'uncleared'),
('862207f1-1d28-40b5-903f-d3dbb312e4b9', '2025-08-14 21:53:51.988009', '2025-08-14 21:53:51.988009', '2025-08-21', '7625b655-732d-49e0-a86b-43994ea89359', NULL, 'Grocery', NULL, 'uncleared'),
('e6469f9d-7388-4ea0-b6f3-e1b085ae1f7f', '2025-08-14 21:50:45.43635', '2025-08-14 21:50:45.43635', '2025-08-14', '7625b655-732d-49e0-a86b-43994ea89359', NULL, 'Transfer to credit card', NULL, 'uncleared');

INSERT INTO "public"."splits" ("id", "created_at", "updated_at", "exchange_id", "destination_register_id", "amount", "counterpart_amount", "memo", "status") VALUES
('118dd813-e2c0-4eed-8dc4-773b87de9b93', '2025-08-14 21:54:36.523784', '2025-08-14 21:54:36.523784', '4c703f3b-7505-4785-9c63-14b38b9e1129', 'af903dfd-3e65-41d2-83f8-db1422b5159c', 1000, -1000, NULL, 'uncleared'),
('2fce88f8-2977-4cb6-89ae-af21ff108e35', '2025-08-14 21:51:40.606261', '2025-08-14 21:51:40.606261', 'e6469f9d-7388-4ea0-b6f3-e1b085ae1f7f', 'f7baf4ac-52f8-494c-87b3-5178b416411f', 12356, -12356, NULL, 'uncleared'),
('34f1be6c-9595-4aa4-a508-0ef1228a7536', '2025-08-14 21:54:36.523784', '2025-08-14 21:54:36.523784', '862207f1-1d28-40b5-903f-d3dbb312e4b9', 'f01ac27a-89f4-4aab-b534-0cf77cee659c', 1000, -1000, NULL, 'uncleared'),
('61dad97b-f7a0-45eb-b325-e6fa5f601937', '2025-08-14 21:54:36.523784', '2025-08-14 21:54:36.523784', '4c703f3b-7505-4785-9c63-14b38b9e1129', 'f01ac27a-89f4-4aab-b534-0cf77cee659c', 10000, -10000, NULL, 'uncleared'),
('88c2dfe2-5589-4695-bd6f-2ccd24bda109', '2025-08-14 21:54:36.523784', '2025-08-14 21:54:36.523784', '4c703f3b-7505-4785-9c63-14b38b9e1129', 'df622be3-dff4-4541-b2c7-40947bc68d5f', 5012, -5012, NULL, 'uncleared'),
('9b861a0a-3637-471e-b24c-d121f5a273ca', '2025-08-14 21:54:36.523784', '2025-08-14 21:54:36.523784', '862207f1-1d28-40b5-903f-d3dbb312e4b9', 'df622be3-dff4-4541-b2c7-40947bc68d5f', 2000, -2000, NULL, 'uncleared'),
('e4261caa-4924-491d-9b3c-2be680f83859', '2025-08-14 21:54:36.523784', '2025-08-14 21:54:36.523784', '862207f1-1d28-40b5-903f-d3dbb312e4b9', 'af903dfd-3e65-41d2-83f8-db1422b5159c', 3000, -3000, NULL, 'uncleared');
