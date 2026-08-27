// Periodic table data — [atomicNumber, symbol, name, category]
// Categories: alkali, alkaline, transition, post-transition, metalloid, nonmetal, halogen, noble, lanthanide, actinide, unknown

export type ElementCategory =
	| 'alkali' | 'alkaline' | 'transition' | 'post-transition' | 'metalloid'
	| 'nonmetal' | 'halogen' | 'noble' | 'lanthanide' | 'actinide' | 'unknown';

export type Element = { z: number; sym: string; name: string; cat: ElementCategory };

export const ELEMENTS: Element[] = [
	{ z: 1, sym: 'H', name: 'Hydrogen', cat: 'nonmetal' },
	{ z: 2, sym: 'He', name: 'Helium', cat: 'noble' },
	{ z: 3, sym: 'Li', name: 'Lithium', cat: 'alkali' },
	{ z: 4, sym: 'Be', name: 'Beryllium', cat: 'alkaline' },
	{ z: 5, sym: 'B', name: 'Boron', cat: 'metalloid' },
	{ z: 6, sym: 'C', name: 'Carbon', cat: 'nonmetal' },
	{ z: 7, sym: 'N', name: 'Nitrogen', cat: 'nonmetal' },
	{ z: 8, sym: 'O', name: 'Oxygen', cat: 'nonmetal' },
	{ z: 9, sym: 'F', name: 'Fluorine', cat: 'halogen' },
	{ z: 10, sym: 'Ne', name: 'Neon', cat: 'noble' },
	{ z: 11, sym: 'Na', name: 'Sodium', cat: 'alkali' },
	{ z: 12, sym: 'Mg', name: 'Magnesium', cat: 'alkaline' },
	{ z: 13, sym: 'Al', name: 'Aluminium', cat: 'post-transition' },
	{ z: 14, sym: 'Si', name: 'Silicon', cat: 'metalloid' },
	{ z: 15, sym: 'P', name: 'Phosphorus', cat: 'nonmetal' },
	{ z: 16, sym: 'S', name: 'Sulfur', cat: 'nonmetal' },
	{ z: 17, sym: 'Cl', name: 'Chlorine', cat: 'halogen' },
	{ z: 18, sym: 'Ar', name: 'Argon', cat: 'noble' },
	{ z: 19, sym: 'K', name: 'Potassium', cat: 'alkali' },
	{ z: 20, sym: 'Ca', name: 'Calcium', cat: 'alkaline' },
	{ z: 21, sym: 'Sc', name: 'Scandium', cat: 'transition' },
	{ z: 22, sym: 'Ti', name: 'Titanium', cat: 'transition' },
	{ z: 23, sym: 'V', name: 'Vanadium', cat: 'transition' },
	{ z: 24, sym: 'Cr', name: 'Chromium', cat: 'transition' },
	{ z: 25, sym: 'Mn', name: 'Manganese', cat: 'transition' },
	{ z: 26, sym: 'Fe', name: 'Iron', cat: 'transition' },
	{ z: 27, sym: 'Co', name: 'Cobalt', cat: 'transition' },
	{ z: 28, sym: 'Ni', name: 'Nickel', cat: 'transition' },
	{ z: 29, sym: 'Cu', name: 'Copper', cat: 'transition' },
	{ z: 30, sym: 'Zn', name: 'Zinc', cat: 'transition' },
	{ z: 31, sym: 'Ga', name: 'Gallium', cat: 'post-transition' },
	{ z: 32, sym: 'Ge', name: 'Germanium', cat: 'metalloid' },
	{ z: 33, sym: 'As', name: 'Arsenic', cat: 'metalloid' },
	{ z: 34, sym: 'Se', name: 'Selenium', cat: 'nonmetal' },
	{ z: 35, sym: 'Br', name: 'Bromine', cat: 'halogen' },
	{ z: 36, sym: 'Kr', name: 'Krypton', cat: 'noble' },
	{ z: 37, sym: 'Rb', name: 'Rubidium', cat: 'alkali' },
	{ z: 38, sym: 'Sr', name: 'Strontium', cat: 'alkaline' },
	{ z: 39, sym: 'Y', name: 'Yttrium', cat: 'transition' },
	{ z: 40, sym: 'Zr', name: 'Zirconium', cat: 'transition' },
	{ z: 41, sym: 'Nb', name: 'Niobium', cat: 'transition' },
	{ z: 42, sym: 'Mo', name: 'Molybdenum', cat: 'transition' },
	{ z: 43, sym: 'Tc', name: 'Technetium', cat: 'transition' },
	{ z: 44, sym: 'Ru', name: 'Ruthenium', cat: 'transition' },
	{ z: 45, sym: 'Rh', name: 'Rhodium', cat: 'transition' },
	{ z: 46, sym: 'Pd', name: 'Palladium', cat: 'transition' },
	{ z: 47, sym: 'Ag', name: 'Silver', cat: 'transition' },
	{ z: 48, sym: 'Cd', name: 'Cadmium', cat: 'transition' },
	{ z: 49, sym: 'In', name: 'Indium', cat: 'post-transition' },
	{ z: 50, sym: 'Sn', name: 'Tin', cat: 'post-transition' },
	{ z: 51, sym: 'Sb', name: 'Antimony', cat: 'metalloid' },
	{ z: 52, sym: 'Te', name: 'Tellurium', cat: 'metalloid' },
	{ z: 53, sym: 'I', name: 'Iodine', cat: 'halogen' },
	{ z: 54, sym: 'Xe', name: 'Xenon', cat: 'noble' },
	{ z: 55, sym: 'Cs', name: 'Caesium', cat: 'alkali' },
	{ z: 56, sym: 'Ba', name: 'Barium', cat: 'alkaline' },
	{ z: 57, sym: 'La', name: 'Lanthanum', cat: 'lanthanide' },
	{ z: 58, sym: 'Ce', name: 'Cerium', cat: 'lanthanide' },
	{ z: 59, sym: 'Pr', name: 'Praseodymium', cat: 'lanthanide' },
	{ z: 60, sym: 'Nd', name: 'Neodymium', cat: 'lanthanide' },
	{ z: 61, sym: 'Pm', name: 'Promethium', cat: 'lanthanide' },
	{ z: 62, sym: 'Sm', name: 'Samarium', cat: 'lanthanide' },
	{ z: 63, sym: 'Eu', name: 'Europium', cat: 'lanthanide' },
	{ z: 64, sym: 'Gd', name: 'Gadolinium', cat: 'lanthanide' },
	{ z: 65, sym: 'Tb', name: 'Terbium', cat: 'lanthanide' },
	{ z: 66, sym: 'Dy', name: 'Dysprosium', cat: 'lanthanide' },
	{ z: 67, sym: 'Ho', name: 'Holmium', cat: 'lanthanide' },
	{ z: 68, sym: 'Er', name: 'Erbium', cat: 'lanthanide' },
	{ z: 69, sym: 'Tm', name: 'Thulium', cat: 'lanthanide' },
	{ z: 70, sym: 'Yb', name: 'Ytterbium', cat: 'lanthanide' },
	{ z: 71, sym: 'Lu', name: 'Lutetium', cat: 'lanthanide' },
	{ z: 72, sym: 'Hf', name: 'Hafnium', cat: 'transition' },
	{ z: 73, sym: 'Ta', name: 'Tantalum', cat: 'transition' },
	{ z: 74, sym: 'W', name: 'Tungsten', cat: 'transition' },
	{ z: 75, sym: 'Re', name: 'Rhenium', cat: 'transition' },
	{ z: 76, sym: 'Os', name: 'Osmium', cat: 'transition' },
	{ z: 77, sym: 'Ir', name: 'Iridium', cat: 'transition' },
	{ z: 78, sym: 'Pt', name: 'Platinum', cat: 'transition' },
	{ z: 79, sym: 'Au', name: 'Gold', cat: 'transition' },
	{ z: 80, sym: 'Hg', name: 'Mercury', cat: 'transition' },
	{ z: 81, sym: 'Tl', name: 'Thallium', cat: 'post-transition' },
	{ z: 82, sym: 'Pb', name: 'Lead', cat: 'post-transition' },
	{ z: 83, sym: 'Bi', name: 'Bismuth', cat: 'post-transition' },
	{ z: 84, sym: 'Po', name: 'Polonium', cat: 'post-transition' },
	{ z: 85, sym: 'At', name: 'Astatine', cat: 'halogen' },
	{ z: 86, sym: 'Rn', name: 'Radon', cat: 'noble' },
	{ z: 87, sym: 'Fr', name: 'Francium', cat: 'alkali' },
	{ z: 88, sym: 'Ra', name: 'Radium', cat: 'alkaline' },
	{ z: 89, sym: 'Ac', name: 'Actinium', cat: 'actinide' },
	{ z: 90, sym: 'Th', name: 'Thorium', cat: 'actinide' },
	{ z: 91, sym: 'Pa', name: 'Protactinium', cat: 'actinide' },
	{ z: 92, sym: 'U', name: 'Uranium', cat: 'actinide' },
	{ z: 93, sym: 'Np', name: 'Neptunium', cat: 'actinide' },
	{ z: 94, sym: 'Pu', name: 'Plutonium', cat: 'actinide' },
	{ z: 95, sym: 'Am', name: 'Americium', cat: 'actinide' },
	{ z: 96, sym: 'Cm', name: 'Curium', cat: 'actinide' },
	{ z: 97, sym: 'Bk', name: 'Berkelium', cat: 'actinide' },
	{ z: 98, sym: 'Cf', name: 'Californium', cat: 'actinide' },
	{ z: 99, sym: 'Es', name: 'Einsteinium', cat: 'actinide' },
	{ z: 100, sym: 'Fm', name: 'Fermium', cat: 'actinide' },
	{ z: 101, sym: 'Md', name: 'Mendelevium', cat: 'actinide' },
	{ z: 102, sym: 'No', name: 'Nobelium', cat: 'actinide' },
	{ z: 103, sym: 'Lr', name: 'Lawrencium', cat: 'actinide' },
	{ z: 104, sym: 'Rf', name: 'Rutherfordium', cat: 'transition' },
	{ z: 105, sym: 'Db', name: 'Dubnium', cat: 'transition' },
	{ z: 106, sym: 'Sg', name: 'Seaborgium', cat: 'transition' },
	{ z: 107, sym: 'Bh', name: 'Bohrium', cat: 'transition' },
	{ z: 108, sym: 'Hs', name: 'Hassium', cat: 'transition' },
	{ z: 109, sym: 'Mt', name: 'Meitnerium', cat: 'unknown' },
	{ z: 110, sym: 'Ds', name: 'Darmstadtium', cat: 'unknown' },
	{ z: 111, sym: 'Rg', name: 'Roentgenium', cat: 'unknown' },
	{ z: 112, sym: 'Cn', name: 'Copernicium', cat: 'unknown' },
	{ z: 113, sym: 'Nh', name: 'Nihonium', cat: 'unknown' },
	{ z: 114, sym: 'Fl', name: 'Flerovium', cat: 'unknown' },
	{ z: 115, sym: 'Mc', name: 'Moscovium', cat: 'unknown' },
	{ z: 116, sym: 'Lv', name: 'Livermorium', cat: 'unknown' },
	{ z: 117, sym: 'Ts', name: 'Tennessine', cat: 'unknown' },
	{ z: 118, sym: 'Og', name: 'Oganesson', cat: 'unknown' },
];

export const CATEGORY_LABELS: Record<ElementCategory, string> = {
	'alkali': 'Alkali metal',
	'alkaline': 'Alkaline earth metal',
	'transition': 'Transition metal',
	'post-transition': 'Post-transition metal',
	'metalloid': 'Metalloid',
	'nonmetal': 'Non-metal',
	'halogen': 'Halogen',
	'noble': 'Noble gas',
	'lanthanide': 'Lanthanide',
	'actinide': 'Actinide',
	'unknown': 'Synthetic',
};

/** Simplified Metal / Non-metal / Metalloid classification. */
export function elementClass(cat: ElementCategory): 'Metal' | 'Non-metal' | 'Metalloid' {
	if (cat === 'metalloid') return 'Metalloid';
	if (cat === 'nonmetal' || cat === 'halogen' || cat === 'noble') return 'Non-metal';
	return 'Metal';
}

export function elementPeriod(z: number): number {
	if (z <= 2) return 1;
	if (z <= 10) return 2;
	if (z <= 18) return 3;
	if (z <= 36) return 4;
	if (z <= 54) return 5;
	if (z <= 86) return 6;
	return 7;
}

export const PERIODIC_DIFFICULTIES = [
	{ id: 30, label: 'Easy', range: 'Elements 1–30', desc: 'Symbols & atomic numbers', tint: 'mint' },
	{ id: 60, label: 'Medium', range: 'Elements 1–60', desc: 'Adds metals & classification', tint: 'amber' },
	{ id: 118, label: 'Hard', range: 'All 118 elements', desc: 'Families, periods & more', tint: 'coral' },
] as const;
