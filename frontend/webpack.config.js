const path = require('path');

module.exports = {
  resolve: {
    extensions: ['.js', '.jsx', '.ts', '.tsx', '.jpg', '.jpeg', '.png'],
  },
  module: {
    rules: [
      {
        test: /\.(png|jpg|jpeg)$/i,
        type: 'asset/resource',
        generator: {
          filename: 'assets/[name][ext][query]',
        },
      },
    ],
  },
};
