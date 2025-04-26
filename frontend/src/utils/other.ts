export function createUrl(host: string, params: string[]) {
  return host + '/' + params.join('/');
}
export function getContrastingColorHex() {
  const randomBrightColor = () => {
    // Генерируем значения от 150 до 255 для создания ярких цветов
    return Math.floor(Math.random() * 106 + 150);
  };

  const r = randomBrightColor();
  const g = randomBrightColor();
  const b = randomBrightColor();

  // Преобразуем RGB в HEX
  const toHex = (value: number) => {
    const hex = value.toString(16);
    return hex.length === 1 ? '0' + hex : hex; // Добавляем ведущий ноль, если нужно
  };

  return '#' + toHex(r) + toHex(g) + toHex(b); // Возвращаем цвет в формате HEX
}

export function createColors() {
  const arr = new Array(50);
  let colors = [];
  for (let e of arr) {
    colors.push(getContrastingColorHex());
  }
  return colors;
}
