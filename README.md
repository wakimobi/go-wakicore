# SERVER

sudo nano /etc/systemd/system/pass_tsel_server.service
sudo nano /etc/systemd/system/pass_tsel_consumer_mo@.service
sudo nano /etc/systemd/system/pass_tsel_consumer_dr@.service
sudo nano /etc/systemd/system/pass_tsel_publisher_renewal.service
sudo nano /etc/systemd/system/pass_tsel_publisher_retry.service
sudo nano /etc/systemd/system/pass_tsel_consumer_renewal@.service
sudo nano /etc/systemd/system/pass_tsel_consumer_retry@.service

=======================================================

# Server

[Unit]
Description=go-pass-tsel server

[Service]
Type=simple
Restart=always
RestartSec=5s
WorkingDirectory=/app/go-pass-tsel
ExecStart=/app/go-pass-tsel/go-pass-tsel server

[Install]
WantedBy=multi-user.target
=======================================================

# Publisher RENEWAL

[Unit]
Description=go-pass-tsel publisher_renewal

[Service]
Type=simple
Restart=always
RestartSec=5s
WorkingDirectory=/app/go-pass-tsel
ExecStart=/app/go-pass-tsel/go-pass-tsel publisher_renewal

[Install]
WantedBy=multi-user.target
=======================================================

# Consumer RENEWAL

[Unit]
Description=go-pass-tsel consumer_renewal

[Service]
Type=simple
Restart=always
RestartSec=5s
WorkingDirectory=/app/go-pass-tsel
ExecStart=/app/go-pass-tsel/go-pass-tsel consumer_renewal

[Install]
WantedBy=multi-user.target
=======================================================

[Unit]
Description=go-pass-tsel consumer_mo %i

[Service]
Type=simple
Restart=always
RestartSec=5s
WorkingDirectory=/app/go-pass-tsel
ExecStart=/app/go-pass-tsel/go-pass-tsel consumer_mo %i

[Install]
WantedBy=multi-user.target
=======================================================

[Unit]
Description=go-pass-tsel consumer_dr %i

[Service]
Type=simple
Restart=always
RestartSec=5s
WorkingDirectory=/app/go-pass-tsel
ExecStart=/app/go-pass-tsel/go-pass-tsel consumer_dr %i

[Install]
WantedBy=multi-user.target
