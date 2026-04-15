BEGIN;

TRUNCATE TABLE cars RESTART IDENTITY;

INSERT INTO cars (make, model, year, description, image_url, gearbox, drivetrain, horsepower, fuel, price_new) VALUES
   ('BMW', 'M3', 2022, 'Performance sedan', 'https://www.netcarshow.com/BMW-M3_Competition_Sedan_M_xDrive-2022-1600-01.jpg', 'automatic', 'RWD', 510, 'gas', 85000),
   ('BMW', 'X5', 2023, 'Luxury SUV', 'https://www.netcarshow.com/BMW-X5-2023-1600-01.jpg', 'automatic', 'AWD', 335, 'hybrid', 78000),
   ('Audi', 'RS5', 2021, 'Sport coupe', 'https://www.netcarshow.com/Audi-RS5_Coupe-2021-1600-01.jpg', 'automatic', 'AWD', 450, 'gas', 75000),
   ('Audi', 'Q7', 2022, 'Large SUV', 'https://www.netcarshow.com/Audi-Q7-2022-1600-01.jpg', 'automatic', 'AWD', 340, 'diesel', 70000),
   ('Mercedes', 'C300', 2021, 'Luxury sedan', 'https://www.netcarshow.com/Mercedes-Benz-C300-2021-1600-01.jpg', 'automatic', 'RWD', 255, 'gas', 55000),
   ('Mercedes', 'GLE', 2022, 'Mid-size SUV', 'https://www.netcarshow.com/Mercedes-Benz-GLE-2022-1600-01.jpg', 'automatic', 'AWD', 362, 'hybrid', 72000),
   ('Toyota', 'Supra', 2021, 'Sports coupe', 'https://www.netcarshow.com/Toyota-Supra-2021-1600-01.jpg', 'automatic', 'RWD', 382, 'gas', 50000),
   ('Toyota', 'RAV4', 2023, 'Compact SUV', 'https://www.netcarshow.com/Toyota-RAV4-2023-1600-01.jpg', 'automatic', 'AWD', 203, 'hybrid', 35000),
   ('Porsche', '911', 2022, 'Iconic sports car', 'https://www.netcarshow.com/Porsche-911_Carrera-2022-1600-01.jpg', 'automatic', 'RWD', 379, 'gas', 110000),
   ('Porsche', 'Cayenne', 2022, 'Luxury SUV', 'https://www.netcarshow.com/Porsche-Cayenne-2022-1600-01.jpg', 'automatic', 'AWD', 335, 'gas', 90000);

COMMIT;