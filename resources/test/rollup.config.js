import typescript from "@rollup/plugin-typescript";
import resolve from "@rollup/plugin-node-resolve";
import commonjs from "@rollup/plugin-commonjs";
import dynamicImportVariables from "@rollup/plugin-dynamic-import-vars";

export default [
  {
    input: "client/src/index.ts",
    output: {
      file: "dist/client.js",
      format: "iife",
    },
    plugins: [
      resolve(),
      commonjs(),
      typescript({
        tsconfig: "client/tsconfig.json",
      }),
    ],
  },
  {
    input: "server/src/index.ts",
    output: {
      file: "dist/server.js",
      format: "cjs",
    },
    plugins: [
      resolve(),
      commonjs(),
      typescript({
        tsconfig: "server/tsconfig.json",
        include: ["server/src/**/*.ts", "../shared/src/**/*.ts"],
      }),
      dynamicImportVariables({
        include: ["server/src/**"],
      }),
    ],
  },
];
