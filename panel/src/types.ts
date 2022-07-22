export interface IMethodInput {
  name: string;
  kind: string;
  type: string;
  fieldList: IMethodInput[];
}

export interface IMethod {
  controller: string;
  fullPath: string;
  httpMethod: string;
  name: string;
  input: IMethodInput;
}
