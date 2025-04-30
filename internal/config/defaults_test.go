package config

// func TestDefaultsParse(t *testing.T) {
// 	t.Parallel()
//
// 	tests := []struct {
// 		input map[string]any
// 		want  Defaults
// 	}{
// 		{},
// 		{
// 			input: map[string]any{"repo": "x"},
// 			want:  Defaults{Repo: github.Repo{Repo: "x"}},
// 		},
// 		{
// 			input: map[string]any{"repo": "x", "owner": "y"},
// 			want: Defaults{
// 				Repo: github.Repo{
// 					Owner: github.User{Login: "y"},
// 					Repo:  "x",
// 				},
// 			},
// 		},
// 		{
// 			input: map[string]any{"owner": "y"},
// 			want: Defaults{
// 				Repo: github.Repo{
// 					Owner: github.User{Login: "y"},
// 				},
// 			},
// 		},
// 		{
// 			input: map[string]any{"repo": "x/y"},
// 			want: Defaults{
// 				Repo: github.Repo{
// 					Owner: github.User{Login: "x"},
// 					Repo:  "y",
// 				},
// 			},
// 		},
// 	}
//
// 	for _, test := range tests {
// 		t.Run(fmt.Sprintf("%v", test.input), func(t *testing.T) {
// 			t.Parallel()
//
// 			d := &Defaults{}
// 			d.parse(test.input)
//
// 			if !test.want.Equal(d) {
// 				t.Errorf("want %+v, but got %+v", test.want, d)
// 			}
// 		})
// 	}
//
// 	t.Run("overwite existing values", func(t *testing.T) {
// 		t.Parallel()
//
// 		input := map[string]any{
// 			"repo": "a/b",
// 		}
//
// 		want := &Defaults{
// 			Repo: github.Repo{
// 				Owner: github.User{Login: "a"},
// 				Repo:  "b",
// 			},
// 		}
//
// 		d := &Defaults{
// 			Repo: github.Repo{
// 				Owner: github.User{Login: "x"},
// 				Repo:  "y",
// 			},
// 		}
//
// 		d.parse(input)
//
// 		if !d.Equal(want) {
// 			t.Errorf("want %+v, but got %+v", want, d)
// 		}
// 	})
// }
