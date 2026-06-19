INSERT INTO krankenhaus (name, mitarbeiteranzahl, bettenanzahl, email, erzeugt, aktualisiert)
VALUES
    ('Staedtisches Klinikum Karlsruhe', 4200, 1500, 'info@klinikum-karlsruhe.example', NOW(), NOW()),
    ('Marienhospital Stuttgart', 3100, 760, 'kontakt@marienhospital-stuttgart.example', NOW(), NOW()),
    ('Universitaetsklinikum Freiburg', 8500, 1600, 'info@uniklinik-freiburg.example', NOW(), NOW());

INSERT INTO adresse (strasse, hausnummer, plz, ort, krankenhaus_id)
VALUES
    ('Moltkestrasse', '90', '76133', 'Karlsruhe', 1000),
    ('Boheimstrasse', '37', '70199', 'Stuttgart', 1001),
    ('Hugstetter Strasse', '55', '79106', 'Freiburg', 1002);

INSERT INTO fachbereich (name, beschreibung, leitung, anzahlaerzte, krankenhaus_id)
VALUES
    ('Kardiologie', 'Herz- und Kreislauferkrankungen', 'Dr. Weber', 28, 1000),
    ('Notaufnahme', 'Zentrale Notfallversorgung', 'Dr. Schmitt', 34, 1000),
    ('Chirurgie', 'Operative Medizin', 'Dr. Keller', 41, 1001),
    ('Onkologie', 'Tumorbehandlung', 'Dr. Schneider', 36, 1002);
