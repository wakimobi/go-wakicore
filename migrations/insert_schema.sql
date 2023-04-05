INSERT INTO services
(category, code, name, price, program_id, sid, renewal_day, trial_day, url_telco, url_portal, url_callback, url_notif_sub, url_notif_unsub, url_notif_renewal, url_postback)
VALUES
('GOALY', 'GOALY', 'GOALY 2', 2220, 'REGGOALY', 'VASREGGOALY_Subs', 2, 0, 'https://api.digitalcore.telkomsel.com', 'https://tsel.goaly.mobi', 'https://tsel.goaly.mobi/subscription/login', 'https://cms-tsel.goaly.mobi/Subscription/notify', 'https://cms-tsel.goaly.mobi/Subscription/unsub_notify', 'https://cms-tsel.goaly.mobi/Subscription/renewal', 'http://kbtools.net/id-passtisel.php');


INSERT INTO contents
(service_id, name, value, tid)
VALUES
(1, 'FIRSTPUSH', '(2220) Gabung di GOALY & dapat reward Exclusive, klik: https://tsel.goaly.mobi (berlaku tarif internet) PIN: @pin Stop: UNREG GOALY ke 99790 CS:02152922391', '4013'),
(1, 'RENEWAL', '(2220) Gabung di GOALY & dapat reward Exclusive, klik: https://tsel.goaly.mobi (berlaku tarif internet) PIN: @pin Stop: UNREG GOALY ke 99790 CS:02152922391', '4013');

INSERT INTO schedules
(id, name, publish_at, unlocked_at, is_unlocked)
VALUES
(1, 'REMINDER', NOW(), NOW(), false),
(2, 'RENEWAL', NOW(), NOW(), false),
(3, 'RETRY', NOW(), NOW(), false);

INSERT INTO adnets
(name, value)
VALUES
('adn', 'adn1');
