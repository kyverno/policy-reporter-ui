export default (hue: number): string => {
  const saturation = '50%';
  const lightness = '40%';
  return `hsl(${hue}, ${saturation}, ${lightness})`;
};
