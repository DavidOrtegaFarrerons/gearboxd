BEGIN;

TRUNCATE TABLE cars RESTART IDENTITY;

INSERT INTO cars (make, model, year, description, image_url, gearbox, drivetrain, horsepower, fuel, price_new) VALUES
   ('BMW', 'M3', 2022, 'Performance sedan', 'https://source.unsplash.com/featured/?bmw,m3', 'automatic', 'RWD', 510, 'gas', 85000),
   ('BMW', 'X5', 2023, 'Luxury SUV', 'https://source.unsplash.com/featured/?bmw,x5', 'automatic', 'AWD', 335, 'hybrid', 78000),
   ('Audi', 'RS5', 2021, 'Sport coupe', 'https://source.unsplash.com/featured/?audi,rs5', 'automatic', 'AWD', 450, 'gas', 75000),
   ('Audi', 'Q7', 2022, 'Large SUV', 'https://source.unsplash.com/featured/?audi,q7', 'automatic', 'AWD', 340, 'diesel', 70000),
   ('Mercedes', 'C300', 2021, 'Luxury sedan', 'https://source.unsplash.com/featured/?mercedes,c300', 'automatic', 'RWD', 255, 'gas', 52000),
   ('Mercedes', 'GLE', 2023, 'Premium SUV', 'https://source.unsplash.com/featured/?mercedes,gle', 'automatic', 'AWD', 362, 'hybrid', 82000),
   ('Toyota', 'Corolla', 2022, 'Reliable compact', 'https://source.unsplash.com/featured/?toyota,corolla', 'CVT', 'FWD', 140, 'hybrid', 25000),
   ('Toyota', 'RAV4', 2023, 'Popular SUV', 'https://source.unsplash.com/featured/?toyota,rav4', 'automatic', 'AWD', 203, 'hybrid', 35000),
   ('Honda', 'Civic', 2021, 'Compact sedan', 'https://source.unsplash.com/featured/?honda,civic', 'manual', 'FWD', 158, 'gas', 24000),
   ('Honda', 'CR-V', 2022, 'Family SUV', 'https://source.unsplash.com/featured/?honda,crv', 'CVT', 'AWD', 190, 'hybrid', 33000),
   ('Tesla', 'Model 3', 2023, 'Electric sedan', 'https://source.unsplash.com/featured/?tesla,model-3', 'automatic', 'RWD', 283, 'electric', 47000),
   ('Tesla', 'Model Y', 2023, 'Electric SUV', 'https://source.unsplash.com/featured/?tesla,model-y', 'automatic', 'AWD', 384, 'electric', 55000),
   ('Ford', 'Mustang', 2021, 'Muscle car', 'https://source.unsplash.com/featured/?ford,mustang', 'manual', 'RWD', 450, 'gas', 55000),
   ('Volkswagen', 'Golf', 2022, 'Classic hatchback', 'https://source.unsplash.com/featured/?volkswagen,golf', 'manual', 'FWD', 130, 'gas', 26000),
   ('Hyundai', 'Ioniq 5', 2023, 'Electric crossover', 'https://source.unsplash.com/featured/?hyundai,ioniq-5', 'automatic', 'AWD', 320, 'electric', 48000),
   ('Kia', 'Sportage', 2023, 'SUV', 'https://source.unsplash.com/featured/?kia,sportage', 'automatic', 'AWD', 180, 'hybrid', 32000),
   ('Peugeot', '3008', 2022, 'SUV', 'https://source.unsplash.com/featured/?peugeot,3008', 'automatic', 'FWD', 130, 'diesel', 31000),
   ('Renault', 'Clio', 2022, 'Compact car', 'https://source.unsplash.com/featured/?renault,clio', 'manual', 'FWD', 100, 'gas', 17000),
   ('Seat', 'Leon', 2022, 'Hatchback', 'https://source.unsplash.com/featured/?seat,leon', 'manual', 'FWD', 130, 'gas', 23000),
   ('Skoda', 'Octavia', 2022, 'Sedan', 'https://source.unsplash.com/featured/?skoda,octavia', 'manual', 'FWD', 150, 'gas', 28000),
   ('Volvo', 'XC60', 2023, 'Luxury SUV', 'https://source.unsplash.com/featured/?volvo,xc60', 'automatic', 'AWD', 250, 'hybrid', 52000),
   ('Mazda', 'MX-5', 2021, 'Roadster', 'https://source.unsplash.com/featured/?mazda,mx5', 'manual', 'RWD', 181, 'gas', 32000),
   ('Nissan', 'Qashqai', 2022, 'Crossover', 'https://source.unsplash.com/featured/?nissan,qashqai', 'CVT', 'FWD', 140, 'gas', 27000),
   ('Jeep', 'Wrangler', 2021, 'Off-road SUV', 'https://source.unsplash.com/featured/?jeep,wrangler', 'manual', '4WD', 285, 'gas', 48000),
   ('Subaru', 'Outback', 2023, 'AWD wagon', 'https://source.unsplash.com/featured/?subaru,outback', 'CVT', 'AWD', 182, 'gas', 36000),
   ('Porsche', 'Taycan', 2023, 'Electric sports', 'https://source.unsplash.com/featured/?porsche,taycan', 'automatic', 'AWD', 408, 'electric', 90000),
   ('Cupra', 'Formentor', 2023, 'Sport SUV', 'https://source.unsplash.com/featured/?cupra,formentor', 'automatic', 'AWD', 310, 'gas', 45000),
   ('Dacia', 'Duster', 2021, 'Budget SUV', 'https://source.unsplash.com/featured/?dacia,duster', 'manual', '4WD', 115, 'diesel', 18000),
   ('Fiat', '500', 2022, 'City car', 'https://source.unsplash.com/featured/?fiat,500', 'manual', 'FWD', 70, 'gas', 16000),
   ('Mini', 'Cooper S', 2021, 'Sport compact', 'https://source.unsplash.com/featured/?mini,cooper', 'manual', 'FWD', 189, 'gas', 32000),
   ('Land Rover', 'Defender', 2023, 'Off-road SUV', 'https://source.unsplash.com/featured/?land-rover,defender', 'automatic', '4WD', 296, 'diesel', 65000);

COMMIT;