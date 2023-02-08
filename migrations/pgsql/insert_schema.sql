INSERT INTO "products"
(id, code, name, auth_user, auth_pass, price, renewal_day, url_notif_sub, url_notif_unsub, url_notif_renewal, url_postback)
VALUES
(1, 'wecarekabe1k', 'WECARE', 'kabe', 'eUpVSUZ16e', 2000, 1, '-', '-', '-', 'http://kbtools.net/id-kb-h3i.php');


INSERT INTO "contents"
(id, product_id, name, value)
VALUES
(1, 1, 'FIRSTPUSH', 'Trims telah brlangganan WECARE. Akses konten kesehatan di wecareid.co gratis konsul dokter dgn promocode PROMOKB1 Stop:UNREG WECARE ke 99345 CS:02122079760'),
(2, 1, 'RENEWAL', 'Trims prpanjangan WECARE brhasil. Akses info kesehatan di wecareid.co gratis konsul dokter dgn promocode PROMOKB1 Stop:UNREG WECARE ke 99345 CS:02122079760');


INSERT INTO "adnets"
(id, name, value)
VALUES
(1, 'kb4', 'xsp'),
(2, 'kb5', 'clc'),
(3, 'kb6', 'sco'),
(4, 'kb7', 'agm2'),
(5, 'kb8', 'mbv');