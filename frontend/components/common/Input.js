import React from "react";

const Input = ({ type, value, isChecked, text, onChange }) => {
  return (
    <p>
      {text}
      <input type={type} value={value} checked={isChecked} onChange={onChange}>
      </input>
    </p>
  );
}

export default Input
