
(defn run [[list-uri output-file]]
  "Expects text file which words are separated by new lines"
  (->>
   (clojure.string/split (slurp list-uri) #"\n")
   (map clojure.string/lower-case)
   (filter #(and (> 7 (count %)) (< 1 (count %))))
   (clojure.string/join "\n")
   (spit output-file)))

(run *command-line-args*)

(println "Done!")
