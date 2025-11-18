INSERT INTO rastochnoy_writedb33 (key, offsett, value) VALUES
('Encoder_X', 2.0, 0.0),
('Encoder_Z', 6.0, 0.0),
('Encoder_Y', 10.0, 0.0),
('Encoder_W', 14.0, 0.0),
('Encoder_V', 18.0, 0.0),
('Padacha_tezlik', 22.0, 0.0),
('Shpindel_tezlik', 26.0, 0.0)
ON CONFLICT (key) DO NOTHING;
