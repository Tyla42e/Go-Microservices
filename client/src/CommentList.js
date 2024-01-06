import React, { useState, useEffect } from "react";
import axios from "axios";

const CommentList = ({ comments }) => {
  console.log("Comments:" + comments)
  if (comments) {
    const renderedComments = comments.map((comment) => {
      return <li key={comment.id}>{comment.content}</li>;
    });
    return <ul>{comments && renderedComments}</ul>;
  }
  return <p>no Data</p>
};

export default CommentList;
