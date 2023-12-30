export const execOnChange = <T>(n: T, o: T, cb: () => any ) => {
  if (JSON.stringify(n) === JSON.stringify(o)) { return; }

  cb()
}

export const onChange = <T>(cb: () => any) => (n: T, o: T) =>  {
  if (JSON.stringify(n) === JSON.stringify(o)) { return; }

  cb()
}
