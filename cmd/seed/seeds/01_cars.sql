BEGIN;

TRUNCATE TABLE cars RESTART IDENTITY;

INSERT INTO cars (make, model, year, description, image_url, gearbox, drivetrain, horsepower, fuel, price_new) VALUES
    ('BMW', 'M3', 2022, 'Performance sedan', 'https://picsum.photos/seed/bmw-m3-2022/800/600', 'automatic', 'RWD', 510, 'gas', 85000),
    ('BMW', 'X5', 2023, 'Luxury SUV', 'https://picsum.photos/seed/bmw-x5-2023/800/600', 'automatic', 'AWD', 335, 'hybrid', 78000),
    ('Audi', 'RS5', 2021, 'Sport coupe', 'https://picsum.photos/seed/audi-rs5-2021/800/600', 'automatic', 'AWD', 450, 'gas', 75000),
    ('Audi', 'Q7', 2022, 'Large SUV', 'https://picsum.photos/seed/audi-q7-2022/800/600', 'automatic', 'AWD', 340, 'diesel', 70000),
    ('Mercedes', 'C300', 2021, 'Luxury sedan', 'https://picsum.photos/seed/mercedes-c300-2021/800/600', 'automatic', 'RWD', 255, 'gas', 55000),
    ('Mercedes', 'GLE', 2022, 'Mid-size SUV', 'https://picsum.photos/seed/mercedes-gle-2022/800/600', 'automatic', 'AWD', 362, 'hybrid', 72000),
    ('Toyota', 'Supra', 2021, 'Sports coupe', 'https://picsum.photos/seed/toyota-supra-2021/800/600', 'automatic', 'RWD', 382, 'gas', 50000),
    ('Toyota', 'RAV4', 2023, 'Compact SUV', 'https://picsum.photos/seed/toyota-rav4-2023/800/600', 'automatic', 'AWD', 203, 'hybrid', 35000),
    ('Porsche', '911', 2022, 'Iconic sports car', 'https://picsum.photos/seed/porsche-911-2022/800/600', 'automatic', 'RWD', 379, 'gas', 110000),
    ('Porsche', 'Cayenne', 2022, 'Luxury SUV', 'https://picsum.photos/seed/porsche-cayenne-2022/800/600', 'automatic', 'AWD', 335, 'gas', 90000);

COMMIT;