export default () => {
  // eslint-disable-next-line no-restricted-globals
  const subPathParts = location.pathname.split('/');

  const subPath = subPathParts.length > 2 && subPathParts[1] !== 'kyverno-plugin' ? `/${subPathParts[1]}/` : '/';

  return subPath;
};
