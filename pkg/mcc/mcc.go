package mcc

type MCC string

func (m MCC) compare(other MCC) int {
	if m < other {
		return -1
	} else if m == other {
		return 0
	}
	return 1
}

func (m MCC) ToCategory() string {
	mcc := map[MCC]string{
		"0000": "Перевод",
		"5411": "Супермаркеты",
		"5533": "Автоуслуги",
		"5812": "Рестораны",
		"5912": "Аптеки",
	}

	if category, ok := mcc[m]; ok {
		return category
	}
	return "Категория не указана"
}
