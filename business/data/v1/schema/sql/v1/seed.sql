INSERT INTO "public"."aggregator"
(id, name, code, api_key, provider_timeout, active, type, payment_by_tgs, logo_url, updated_at, created_at, deleted_at)
VALUES ('b128d272-3f64-4aa4-bcb3-5909cfb1dd05', 'poc-tgs', 'poc-tgs', '80c084e7-8038-4c91-83a3-9ec8ff4315af', 5000,
        true, 'default', true,
        'https://image.shutterstock.com/image-vector/stripes-dots-circle-logo-vector-260nw-1139400776.jpg',
        TIMESTAMP '2022-06-04 09:47:15.6567194', TIMESTAMP '2022-06-04 09:47:15.6567194',
        null);

INSERT INTO "public"."enterprise_pack"
(id, name, send_monthly_report, can_customize_report, included_field_in_report, send_expense_report,
 max_active_collaborator_per_month, created_at, updated_at,
 deleted_at)
VALUES ('60a82100-72a0-413a-b91a-d7acdf5bc511', 'siege', true, true, 'test', true, 120,
        TIMESTAMP '2022-06-04 09:47:57.5642897', TIMESTAMP '2022-06-04 09:47:57.5642897',
        null);

INSERT INTO "public"."enterprise"
(id, pack_id, code, name, blocked_provider, blocked_product_type, contact_email,
 description, logo_url, max_carbon_emission, active, created_at, updated_at,
 deleted_at)
VALUES ('b8f01d16-8675-4e61-8d04-8e010a40ebec', '60a82100-72a0-413a-b91a-d7acdf5bc511', '331822', 'samaketech', 'bolt',
        'eco', 'opmadou@gmail.com', 'This is a test enterprise',
        'https://image.shutterstock.com/image-vector/stripes-dots-circle-logo-vector-260nw-1139400776.jpg', 5000, true,
        TIMESTAMP '2022-06-04 09:49:24.0355278', TIMESTAMP '2022-06-04 09:49:24.0355278',
        null);

INSERT INTO enterprise_policy
(id, enterprise_id, name, blocked_days, description, collaborator_budget,
 start_service_time, end_service_time, budget_type, created_at, updated_at, deleted_at)
VALUES ('c34061ae-a6d3-4147-a351-e85f511ac15f', 'b8f01d16-8675-4e61-8d04-8e010a40ebec', 'siege', '',
        'all collaborator that are admin', 1000, '07:00:00', '23:00:00', 'monthly',
        TIMESTAMP '2022-06-04 09:50:08.6986978', TIMESTAMP '2022-06-04 09:50:08.6986978',
        null);

INSERT INTO "public"."invoicing_entity"
(id, enterprise_id, denomination, number, vat, street, postal_code, town,
 country, created_at, updated_at, deleted_at)
VALUES ('55d9cb47-23e6-48f1-9321-eefaec9d1fba', 'b8f01d16-8675-4e61-8d04-8e010a40ebec', 'samaketech', 'FE500RE',
        'FEVAT556022', '11 ya parkovaya', '105077', 'paris', 'france',
        TIMESTAMP '2022-06-04 09:50:50.2304067',
        TIMESTAMP '2022-06-04 09:50:50.2304067', null);

INSERT INTO "public"."enterprise_team"
(id, name, invoicing_entity_id, policy_id, enterprise_id, description,
 payment_method, created_at, updated_at, deleted_at)
VALUES ('68b57953-6a12-454e-99b0-a1b6a4bcf58b', 'siege', '55d9cb47-23e6-48f1-9321-eefaec9d1fba',
        'c34061ae-a6d3-4147-a351-e85f511ac15f', 'b8f01d16-8675-4e61-8d04-8e010a40ebec',
        'all collaborator that are admin', 'NDF', TIMESTAMP '2022-06-04 09:51:55.8903145',
        TIMESTAMP '2022-06-04 09:51:55.8903145', null);

INSERT INTO "public"."user"
(id, aggregator_id, email, phone_number, name, stripe_id,
 api_key, active, cognito_id, is_monthly_active, is_cgu_accepted, role, created_at,
 updated_at, deleted_at, budget)
VALUES ('7824b480-6342-4e28-8c48-7011bd051712', 'b128d272-3f64-4aa4-bcb3-5909cfb1dd05', 'opmadou@gmail.com',
        '+33757130289', 'Mahamadou Samake', 'cus_LoWisF53KJsoa5', '80c084e7-8038-4c91-83a3-9ec8ff4315af', true,
        '6286d251a91f5ad486de1a012202b8cbd70a6909963816b6a1b017ea20df8807', false, true, 'USER',
        TIMESTAMP '2022-06-04 09:52:25.9773366', TIMESTAMP '2022-06-04 09:52:25.9773366',
        null, 0);

INSERT INTO public.payment_method (id, name, user_id, display_credit_card, stripe_id, type, is_favorite, created_at,
                                   updated_at, deleted_at)
VALUES ('5658a7ce-37eb-4cda-a35b-c4547129c2c7', 'Mahamadou Samake', '7824b480-6342-4e28-8c48-7011bd051712',
        'XXXX XXXX XXXX 4242', 'pm_1L6v4yCSGvJXXYWCtvFfFtoi', 'Visa', true, '2022-06-04 11:23:40.909482',
        '2022-06-04 11:23:40.909482', null);
