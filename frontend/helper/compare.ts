export const execOnChange = <T>(n: T, o: T, cb: () => any ) => {
  if (JSON.stringify(n) === JSON.stringify(o)) { return; }

  cb()
}
