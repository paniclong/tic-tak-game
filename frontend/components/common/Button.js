import React from "react";

const Button = ({ value, onClick, styles, disabled }) => {
  return (
    <button className="Button" onClick={onClick} style={styles} disabled={disabled}>
      {value}
    </button>
  );
}

export default Button
