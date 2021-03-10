INSERT INTO clients(id, login, password, first_name, last_name, middle_name, passport, birthday, status)
VALUES (1, 'user1', 'hash1', 'Юзер', 'Юзеров', 'Юзерович', '1111000222', '2000-01-01', 'INACTIVE'),
       (2, 'user2', 'hash2', 'Владимир', 'Владимиров', 'Владимирович', '5648111222', '1960-06-25', 'ACTIVE');

INSERT INTO cards(id, number, balance, issuer, holder, owner_id, status)
VALUES (1, '4561 2612 1234 5464', 1000000, 'Visa', 'Владимир Владимирович', 2, 'ACTIVE'),
       (2, '2612 4561 1234 5464', 0, 'MasterCard', 'Юзвер', 1, 'INACTIVE'),
       (3, '2612 4561 1254 3464', 500000, 'MasterCard', 'Uzer Uzerov', 1, 'ACTIVE');

INSERT INTO icons(id, title, uri)
VALUES (1, 'Альфа', 'https://...'),
       (2, 'Продукты', 'https://...'),
       (3, 'Мегафон', 'https://...'),
       (4, 'Тинькофф', 'https://...');

INSERT INTO mccs(id, text)
VALUES ('4814', 'Мобильная связь'),
       ('5411', 'Супермаркеты'),
       ('5533', 'Автоуслуги'),
       ('5812', 'Рестораны'),
       ('5912', 'Аптеки'),
       ('5732', 'Магазины электроники');

INSERT INTO transactions(from_id, to_id, sum, mcc_id, icon_id, description)
VALUES (NULL, 1, 5000000, NULL, 1, 'Пополнение через Альфа-Банк'),
       (1, NULL, -100000, '5411', 2, 'Продукты'),
       (1, NULL, -100000, '4814', 3, 'Пополнение телефона'),
       (1, NULL, -100000, '5411', 2, 'Продукты'),
       (1, NULL, -200000, '5411', 2, 'Продукты'),
       (1, NULL, -150000, '5411', 2, 'Продукты'),
       (1, 3, -1000000, NULL, 4, 'Перевод'),
       (1, NULL, -120000, '5411', 2, 'Продукты');
