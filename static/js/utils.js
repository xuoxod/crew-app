export const stringify = (obj) => {
  if (null != obj) {
    return JSON.stringify(obj);
  }
  return null;
};

export const parse = (jsonStr) => {
  if (null != jsonStr) {
    return JSON.parse(jsonStr);
  }
  return null;
};

export const keys = (obj) => {
  if (null != obj && !Array.isArray(obj) && typeof obj == "object") {
    return Object.keys(obj);
  }
  return null;
};

export const cap = (arg) => {
  let word_split = null,
    line = "",
    word = arg.toString();
  if (null !== word && undefined !== word) {
    if (
      word.trim().toLowerCase() === "id" ||
      word.trim().toLowerCase() === "ssn" ||
      word.trim().toLowerCase() === "sku" ||
      word.trim().toLowerCase() === "vm" ||
      word.trim().toLowerCase() === "mac" ||
      word.trim().toLowerCase() === "imei" ||
      word.trim().toLowerCase() === "os" ||
      word.trim().toLowerCase() === "atm" ||
      word.trim().toLowerCase() === "pa" ||
      word.trim().toLowerCase() === "rjw"
    ) {
      return word.toUpperCase();
    } else if (word.match(/[-]/)) {
      if (null !== (word_split = word.split(["-"])).length > 0) {
        for (let i = 0; i < word_split.length; i++) {
          if (i < word_split.length - 1) {
            line +=
              word_split[i].substring(0, 1).toUpperCase() +
              word_split[i].substring(1) +
              "-";
          } else {
            line +=
              word_split[i].substring(0, 1).toUpperCase() +
              word_split[i].substring(1);
          }
        }
        return line;
      }
    } else if (word.match(/[ ]/)) {
      if (null !== (word_split = word.split([" "])).length > 0) {
        for (let i = 0; i < word_split.length; i++) {
          if (i < word_split.length - 1) {
            line +=
              word_split[i].substring(0, 1).toUpperCase() +
              word_split[i].substring(1) +
              " ";
          } else {
            line +=
              word_split[i].substring(0, 1).toUpperCase() +
              word_split[i].substring(1);
          }
        }
        return line;
      }
    } else {
      return word.substring(0, 1).toUpperCase() + word.substring(1);
    }
  }
};

export const stripTags = (input) => {
  return input.replace(/(<(?:.|\n)*?>|^[^\w])/gm, "");
};

export const size = (arg = null) => {
  if (null != arg) {
    if (Array.isArray(arg)) {
      return arg.length;
    } else if (arg instanceof Object && !Array.isArray(arg)) {
      return Object.keys(arg).length;
    } else if (
      !(arg instanceof Object) &&
      !Array.isArray(arg) &&
      typeof arg == "string"
    ) {
      return arg.length;
    } else {
      return NaN;
    }
  }
};

export const isObject = (data) =>
  typeof data === "object" && !Array.isArray(data) && null != data;

export const isArray = (data) => null != data && Array.isArray(data);

export const isNull = (data) => null == data || undefined == data;
