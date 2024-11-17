//go:build !ignore_autogenerated

/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by controller-gen. DO NOT EDIT.

package v1

import (
	"k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Action) DeepCopyInto(out *Action) {
	*out = *in
	if in.Send != nil {
		in, out := &in.Send, &out.Send
		*out = new(SendAction)
		(*in).DeepCopyInto(*out)
	}
	if in.Concat != nil {
		in, out := &in.Concat, &out.Concat
		*out = new(ConcatAction)
		(*in).DeepCopyInto(*out)
	}
	if in.Remove != nil {
		in, out := &in.Remove, &out.Remove
		*out = new(RemoveAction)
		(*in).DeepCopyInto(*out)
	}
	if in.Rename != nil {
		in, out := &in.Rename, &out.Rename
		*out = new(RenameAction)
		**out = **in
	}
	if in.Duplicate != nil {
		in, out := &in.Duplicate, &out.Duplicate
		*out = new(DuplicateAction)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Action.
func (in *Action) DeepCopy() *Action {
	if in == nil {
		return nil
	}
	out := new(Action)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ConcatAction) DeepCopyInto(out *ConcatAction) {
	*out = *in
	if in.InputKeys != nil {
		in, out := &in.InputKeys, &out.InputKeys
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ConcatAction.
func (in *ConcatAction) DeepCopy() *ConcatAction {
	if in == nil {
		return nil
	}
	out := new(ConcatAction)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Decide) DeepCopyInto(out *Decide) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Decide.
func (in *Decide) DeepCopy() *Decide {
	if in == nil {
		return nil
	}
	out := new(Decide)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Decide) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DecideList) DeepCopyInto(out *DecideList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Decide, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DecideList.
func (in *DecideList) DeepCopy() *DecideList {
	if in == nil {
		return nil
	}
	out := new(DecideList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *DecideList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DecideSpec) DeepCopyInto(out *DecideSpec) {
	*out = *in
	if in.Actions != nil {
		in, out := &in.Actions, &out.Actions
		*out = make([]Action, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Next != nil {
		in, out := &in.Next, &out.Next
		*out = make([]Next, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DecideSpec.
func (in *DecideSpec) DeepCopy() *DecideSpec {
	if in == nil {
		return nil
	}
	out := new(DecideSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DecideStatus) DeepCopyInto(out *DecideStatus) {
	*out = *in
	in.Input.DeepCopyInto(&out.Input)
	in.LastUpdated.DeepCopyInto(&out.LastUpdated)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DecideStatus.
func (in *DecideStatus) DeepCopy() *DecideStatus {
	if in == nil {
		return nil
	}
	out := new(DecideStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Destination) DeepCopyInto(out *Destination) {
	*out = *in
	if in.HTTP != nil {
		in, out := &in.HTTP, &out.HTTP
		*out = new(HTTPDestination)
		**out = **in
	}
	if in.File != nil {
		in, out := &in.File, &out.File
		*out = new(FileDestination)
		**out = **in
	}
	if in.GRPC != nil {
		in, out := &in.GRPC, &out.GRPC
		*out = new(GRPCDestination)
		**out = **in
	}
	if in.Opa != nil {
		in, out := &in.Opa, &out.Opa
		*out = new(OpaDestination)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Destination.
func (in *Destination) DeepCopy() *Destination {
	if in == nil {
		return nil
	}
	out := new(Destination)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DuplicateAction) DeepCopyInto(out *DuplicateAction) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DuplicateAction.
func (in *DuplicateAction) DeepCopy() *DuplicateAction {
	if in == nil {
		return nil
	}
	out := new(DuplicateAction)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Element) DeepCopyInto(out *Element) {
	*out = *in
	if in.Observe != nil {
		in, out := &in.Observe, &out.Observe
		*out = new(ObserveSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.Decide != nil {
		in, out := &in.Decide, &out.Decide
		*out = new(DecideSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.Execute != nil {
		in, out := &in.Execute, &out.Execute
		*out = new(ExecuteSpec)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Element.
func (in *Element) DeepCopy() *Element {
	if in == nil {
		return nil
	}
	out := new(Element)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Execute) DeepCopyInto(out *Execute) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Execute.
func (in *Execute) DeepCopy() *Execute {
	if in == nil {
		return nil
	}
	out := new(Execute)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Execute) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ExecuteList) DeepCopyInto(out *ExecuteList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Execute, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ExecuteList.
func (in *ExecuteList) DeepCopy() *ExecuteList {
	if in == nil {
		return nil
	}
	out := new(ExecuteList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ExecuteList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ExecuteSpec) DeepCopyInto(out *ExecuteSpec) {
	*out = *in
	in.Destination.DeepCopyInto(&out.Destination)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ExecuteSpec.
func (in *ExecuteSpec) DeepCopy() *ExecuteSpec {
	if in == nil {
		return nil
	}
	out := new(ExecuteSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ExecuteStatus) DeepCopyInto(out *ExecuteStatus) {
	*out = *in
	in.Input.DeepCopyInto(&out.Input)
	in.LastUpdated.DeepCopyInto(&out.LastUpdated)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ExecuteStatus.
func (in *ExecuteStatus) DeepCopy() *ExecuteStatus {
	if in == nil {
		return nil
	}
	out := new(ExecuteStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FileDestination) DeepCopyInto(out *FileDestination) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FileDestination.
func (in *FileDestination) DeepCopy() *FileDestination {
	if in == nil {
		return nil
	}
	out := new(FileDestination)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GRPCDestination) DeepCopyInto(out *GRPCDestination) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GRPCDestination.
func (in *GRPCDestination) DeepCopy() *GRPCDestination {
	if in == nil {
		return nil
	}
	out := new(GRPCDestination)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HTTPDestination) DeepCopyInto(out *HTTPDestination) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HTTPDestination.
func (in *HTTPDestination) DeepCopy() *HTTPDestination {
	if in == nil {
		return nil
	}
	out := new(HTTPDestination)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Master) DeepCopyInto(out *Master) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Master.
func (in *Master) DeepCopy() *Master {
	if in == nil {
		return nil
	}
	out := new(Master)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Master) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MasterList) DeepCopyInto(out *MasterList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Master, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MasterList.
func (in *MasterList) DeepCopy() *MasterList {
	if in == nil {
		return nil
	}
	out := new(MasterList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *MasterList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MasterSpec) DeepCopyInto(out *MasterSpec) {
	*out = *in
	if in.Elements != nil {
		in, out := &in.Elements, &out.Elements
		*out = make([]Element, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MasterSpec.
func (in *MasterSpec) DeepCopy() *MasterSpec {
	if in == nil {
		return nil
	}
	out := new(MasterSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MasterStatus) DeepCopyInto(out *MasterStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MasterStatus.
func (in *MasterStatus) DeepCopy() *MasterStatus {
	if in == nil {
		return nil
	}
	out := new(MasterStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Next) DeepCopyInto(out *Next) {
	*out = *in
	if in.Keys != nil {
		in, out := &in.Keys, &out.Keys
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Next.
func (in *Next) DeepCopy() *Next {
	if in == nil {
		return nil
	}
	out := new(Next)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Observe) DeepCopyInto(out *Observe) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Observe.
func (in *Observe) DeepCopy() *Observe {
	if in == nil {
		return nil
	}
	out := new(Observe)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Observe) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ObserveList) DeepCopyInto(out *ObserveList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Observe, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ObserveList.
func (in *ObserveList) DeepCopy() *ObserveList {
	if in == nil {
		return nil
	}
	out := new(ObserveList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ObserveList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ObserveSpec) DeepCopyInto(out *ObserveSpec) {
	*out = *in
	if in.Next != nil {
		in, out := &in.Next, &out.Next
		*out = make([]Next, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ObserveSpec.
func (in *ObserveSpec) DeepCopy() *ObserveSpec {
	if in == nil {
		return nil
	}
	out := new(ObserveSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ObserveStatus) DeepCopyInto(out *ObserveStatus) {
	*out = *in
	in.Input.DeepCopyInto(&out.Input)
	in.LastUpdated.DeepCopyInto(&out.LastUpdated)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ObserveStatus.
func (in *ObserveStatus) DeepCopy() *ObserveStatus {
	if in == nil {
		return nil
	}
	out := new(ObserveStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *OpaDestination) DeepCopyInto(out *OpaDestination) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new OpaDestination.
func (in *OpaDestination) DeepCopy() *OpaDestination {
	if in == nil {
		return nil
	}
	out := new(OpaDestination)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RemoveAction) DeepCopyInto(out *RemoveAction) {
	*out = *in
	if in.InputKeys != nil {
		in, out := &in.InputKeys, &out.InputKeys
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RemoveAction.
func (in *RemoveAction) DeepCopy() *RemoveAction {
	if in == nil {
		return nil
	}
	out := new(RemoveAction)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RenameAction) DeepCopyInto(out *RenameAction) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RenameAction.
func (in *RenameAction) DeepCopy() *RenameAction {
	if in == nil {
		return nil
	}
	out := new(RenameAction)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SendAction) DeepCopyInto(out *SendAction) {
	*out = *in
	in.Destination.DeepCopyInto(&out.Destination)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SendAction.
func (in *SendAction) DeepCopy() *SendAction {
	if in == nil {
		return nil
	}
	out := new(SendAction)
	in.DeepCopyInto(out)
	return out
}
